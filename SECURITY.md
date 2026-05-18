# Security Policy

agent-help output is plain text and may be logged by agent runtimes, shells, terminals, CI systems, or observability tooling.

## Reporting Security Issues

If you find a security issue in the specification, examples, or reference implementation, please report it privately to the maintainers when possible. If no private contact is available, open a GitHub issue with minimal details and ask for a private follow-up path.

Do not include secrets, exploit payloads, credentials, or private data in public issues.

## Guidance for Implementers

Agent-facing output should be treated as potentially visible to an LLM agent and any systems that record the agent session.

Implementations should:

- Redact API keys, tokens, passwords, cookies, SSH keys, and cloud credentials.
- Avoid printing unnecessary personal data.
- Avoid embedding raw environment variables in `--agent-out`.
- Avoid returning hidden files or secret paths unless explicitly requested.
- Make pagination and truncation explicit so agents do not over-request sensitive data.
- Prefer stable opaque IDs over sensitive internal identifiers.

`--agent-help` should describe how to use the CLI, not expose runtime secrets or local machine details.

