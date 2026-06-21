package permission

import (
	"path"
	"slices"
	"strings"
	"unicode"

	toolruntime "forgecode/internal/tool-runtime"
)

type bashAnalyzer struct{}

type bashToken struct {
	Text     string
	Operator bool
	Quoted   bool
}

type bashCommand struct {
	Program   string
	Args      []bashToken
	SepBefore string
}

func NewBashAnalyzer() BashAnalyzer {
	return bashAnalyzer{}
}

func (bashAnalyzer) Analyze(cmd string) (BashAnalysis, error) {
	if strings.TrimSpace(cmd) == "" {
		return BashAnalysis{}, toolruntime.NewError(toolruntime.ValidationError, "bash command is required")
	}

	tokens, analysis := lexBash(cmd)
	commands := parseBashCommands(tokens, &analysis)
	classifyBashCommands(commands, &analysis)
	return analysis, nil
}

func lexBash(cmd string) ([]bashToken, BashAnalysis) {
	var tokens []bashToken
	var analysis BashAnalysis
	var buf strings.Builder
	var quoted, single, double, escape bool
	runes := []rune(cmd)

	flush := func() {
		if buf.Len() == 0 {
			quoted = false
			return
		}
		tokens = append(tokens, bashToken{Text: buf.String(), Quoted: quoted})
		buf.Reset()
		quoted = false
	}
	addOperator := func(text string) {
		flush()
		tokens = append(tokens, bashToken{Text: text, Operator: true})
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if escape {
			buf.WriteRune(r)
			escape = false
			continue
		}
		if single {
			if r == '\'' {
				single = false
				continue
			}
			buf.WriteRune(r)
			continue
		}
		if double {
			switch r {
			case '\\':
				escape = true
			case '"':
				double = false
			case '$':
				if i+1 < len(runes) && runes[i+1] == '(' {
					analysis.CmdSubst = true
					analysis.Subshells = true
					addOperator("$(")
					i++
					continue
				}
				buf.WriteRune(r)
			case '`':
				analysis.CmdSubst = true
				addOperator("`")
			default:
				buf.WriteRune(r)
			}
			continue
		}

		switch {
		case unicode.IsSpace(r):
			flush()
		case r == '\\':
			escape = true
		case r == '\'':
			quoted = true
			single = true
		case r == '"':
			quoted = true
			double = true
		case r == '$' && i+1 < len(runes) && runes[i+1] == '(':
			analysis.CmdSubst = true
			analysis.Subshells = true
			addOperator("$(")
			i++
		case r == '`':
			analysis.CmdSubst = true
			addOperator("`")
		case r == '(' || r == ')':
			analysis.Subshells = true
			addOperator(string(r))
		case r == '|':
			analysis.Pipes = true
			if i+1 < len(runes) && runes[i+1] == '|' {
				addOperator("||")
				i++
			} else if i+1 < len(runes) && runes[i+1] == '&' {
				addOperator("|&")
				i++
			} else {
				addOperator("|")
			}
		case r == '&':
			if i+1 < len(runes) && runes[i+1] == '&' {
				addOperator("&&")
				i++
			} else {
				addOperator("&")
			}
		case r == ';':
			addOperator(";")
		case r == '<' || r == '>':
			analysis.Redirects = true
			op := string(r)
			if i+1 < len(runes) && runes[i+1] == r {
				op += string(r)
				i++
			}
			addOperator(op)
		default:
			buf.WriteRune(r)
		}
	}
	flush()
	return tokens, analysis
}

func parseBashCommands(tokens []bashToken, analysis *BashAnalysis) []bashCommand {
	var commands []bashCommand
	var current *bashCommand
	commandStart := true
	redirectTarget := false
	sepBefore := ""

	for _, token := range tokens {
		if token.Operator {
			switch token.Text {
			case "|", "|&", "||", "&&", ";", "&", "$(", "(", "`":
				commandStart = true
				current = nil
				sepBefore = token.Text
			case ")", ">", ">>", "<", "<<":
				if isRedirectOperator(token.Text) {
					redirectTarget = true
				}
			default:
				if isRedirectOperator(token.Text) {
					redirectTarget = true
				}
			}
			continue
		}
		if redirectTarget {
			redirectTarget = false
			continue
		}
		if commandStart {
			if isAssignment(token.Text) {
				continue
			}
			program := normalizeProgram(token.Text)
			if program == "" {
				continue
			}
			commands = append(commands, bashCommand{Program: program, SepBefore: sepBefore})
			if !slices.Contains(analysis.Programs, program) {
				analysis.Programs = append(analysis.Programs, program)
			}
			current = &commands[len(commands)-1]
			commandStart = false
			sepBefore = ""
			continue
		}
		if current != nil {
			current.Args = append(current.Args, token)
		}
	}
	return commands
}

func classifyBashCommands(commands []bashCommand, analysis *BashAnalysis) {
	networkCommandSeenBefore := false
	for i := range commands {
		command := commands[i]
		args := lowerArgs(command.Args)
		commandHasNetwork := isNetworkProgram(command.Program) || hasUnquotedURL(command.Args) || isPowerShellNetwork(command)

		if commandHasNetwork {
			analysis.NetworkAccess = true
		}
		if isDeletionCommand(command.Program, args) || isWrappedDeletion(command.Program, args) {
			analysis.FileDeletion = true
		}
		if isForcePush(command.Program, args) {
			analysis.ForcePush = true
		}
		if isDockerProgram(command.Program) {
			analysis.Docker = true
		}
		if isKubernetesProgram(command.Program) {
			analysis.Kubernetes = true
		}
		if isDBWrite(command.Program, args) {
			analysis.DBWrite = true
		}
		if isPrivilegeEscalation(command.Program, args) {
			analysis.PrivilegeEscalation = true
		}
		if networkCommandSeenBefore && isExecutorProgram(command.Program) {
			analysis.DownloadThenExec = true
		}
		if i > 0 && isNetworkProgram(commands[i-1].Program) && isExecutorProgram(command.Program) {
			analysis.DownloadThenExec = true
		}
		if commandHasNetwork && isPowerShellExpressionExec(command) {
			analysis.DownloadThenExec = true
		}
		if commandHasNetwork {
			networkCommandSeenBefore = true
		}
	}
}

func normalizeProgram(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.Trim(cleaned, `"'`)
	cleaned = strings.ReplaceAll(cleaned, `\`, "/")
	cleaned = strings.ToLower(path.Base(cleaned))
	for _, suffix := range []string{".exe", ".cmd", ".bat", ".ps1"} {
		cleaned = strings.TrimSuffix(cleaned, suffix)
	}
	return cleaned
}

func isAssignment(text string) bool {
	eq := strings.IndexRune(text, '=')
	if eq <= 0 {
		return false
	}
	for _, r := range text[:eq] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func isRedirectOperator(text string) bool {
	return strings.Contains(text, ">") || strings.Contains(text, "<")
}

func lowerArgs(args []bashToken) []string {
	lowered := make([]string, 0, len(args))
	for _, arg := range args {
		lowered = append(lowered, strings.ToLower(arg.Text))
	}
	return lowered
}

func hasUnquotedURL(args []bashToken) bool {
	for _, arg := range args {
		if arg.Quoted {
			continue
		}
		lower := strings.ToLower(arg.Text)
		if strings.Contains(lower, "http://") || strings.Contains(lower, "https://") {
			return true
		}
	}
	return false
}

func isNetworkProgram(program string) bool {
	switch program {
	case "curl", "wget", "ssh", "scp", "sftp", "ftp", "nc", "ncat", "telnet", "iwr", "irm", "invoke-webrequest", "invoke-restmethod":
		return true
	default:
		return false
	}
}

func isDeletionCommand(program string, args []string) bool {
	switch program {
	case "rm", "del", "erase", "rmdir", "remove-item", "ri":
		return true
	case "find":
		return slices.Contains(args, "-delete")
	default:
		return false
	}
}

func isWrappedDeletion(program string, args []string) bool {
	if program != "sudo" && program != "runas" {
		return false
	}
	for i, arg := range args {
		nested := normalizeProgram(arg)
		if nested == "" || strings.HasPrefix(nested, "-") {
			continue
		}
		return isDeletionCommand(nested, args[i+1:])
	}
	return false
}

func isForcePush(program string, args []string) bool {
	if program != "git" || !slices.Contains(args, "push") {
		return false
	}
	for _, arg := range args {
		if arg == "-f" || arg == "--force" || arg == "--force-with-lease" || strings.HasPrefix(arg, "+") {
			return true
		}
	}
	return false
}

func isDockerProgram(program string) bool {
	switch program {
	case "docker", "docker-compose", "podman":
		return true
	default:
		return false
	}
}

func isKubernetesProgram(program string) bool {
	switch program {
	case "kubectl", "helm", "kustomize":
		return true
	default:
		return false
	}
}

func isDBWrite(program string, args []string) bool {
	switch program {
	case "psql", "mysql", "sqlite3", "mongosh", "mongo":
	default:
		return false
	}
	joined := " " + strings.Join(args, " ") + " "
	for _, keyword := range []string{"insert", "update", "delete", "drop", "alter", "create", "truncate", "grant", "revoke"} {
		if strings.Contains(joined, " "+keyword+" ") || strings.Contains(joined, keyword+" ") {
			return true
		}
	}
	return false
}

func isPrivilegeEscalation(program string, args []string) bool {
	switch program {
	case "sudo", "su", "runas":
		return true
	case "powershell", "pwsh", "start-process":
		joined := strings.Join(args, " ")
		return strings.Contains(joined, "runas")
	default:
		return false
	}
}

func isExecutorProgram(program string) bool {
	switch program {
	case "sh", "bash", "zsh", "fish", "cmd", "powershell", "pwsh", "iex", "invoke-expression", "python", "python3", "node", "ruby", "perl":
		return true
	default:
		return false
	}
}

func isPowerShellNetwork(command bashCommand) bool {
	if command.Program != "powershell" && command.Program != "pwsh" {
		return false
	}
	joined := strings.ToLower(joinTokenText(command.Args))
	return strings.Contains(joined, "invoke-webrequest") ||
		strings.Contains(joined, "invoke-restmethod") ||
		strings.Contains(joined, " iwr ") ||
		strings.Contains(joined, " irm ") ||
		strings.Contains(joined, "http://") ||
		strings.Contains(joined, "https://")
}

func isPowerShellExpressionExec(command bashCommand) bool {
	if command.Program != "powershell" && command.Program != "pwsh" {
		return false
	}
	joined := " " + strings.ToLower(joinTokenText(command.Args)) + " "
	return strings.Contains(joined, " iex ") || strings.Contains(joined, "invoke-expression")
}

func joinTokenText(tokens []bashToken) string {
	parts := make([]string, 0, len(tokens))
	for _, token := range tokens {
		parts = append(parts, token.Text)
	}
	return strings.Join(parts, " ")
}
