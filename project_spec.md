# Project: Guardrail (gr)

**Objective:** Build a context-aware CLI tool in Go to safely manage environment switching and prevent accidental destructive commands across several distinct cluster environments.

## Tech Stack & Libraries

- **Language:** Go 1.24+

- **CLI Framework:** Cobra (`spf13/cobra`)
- **Config Management:** Viper (`spf13/viper`) for YAML parsing
- **UI/Terminal:** `fatih/color` for visual environment indicators, `manifoldco/promptui` for interactive safety prompts.

## Core Requirements

1. **Config Driven:** Must read from `~/.guardrail/config.yaml` which defines Cloud profiles, Kubernetes contexts, and Terraform workspaces for the several clusters.

2. **Context Switching:** The `gr use <env>` command must synchronously update the local shell environment, KUBECONFIG, and Cloud context.

3. **Safety Interceptor:** Must wrap execution commands. If an environment is flagged as `is_production: true`, commands containing `delete`, `uninstall`, `destroy`, or `apply` must trigger a manual string-match confirmation prompt.

## Developer Context

The primary developer is transitioning from Python to Go. All code generated must include brief, inline comments explaining Go-specific idioms (e.g., pointers, interfaces, error handling, goroutines) and how they compare to Python equivalents.
