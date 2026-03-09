# Guardrail 🛡️

Guardrail is a context-aware CLI tool written in Go that safely manages environment switching across various clusters and cloud accounts. Its core feature is a **Safety Interceptor** that halts execution when running high-risk commands (like `delete`, `uninstall`, or `destroy`) against a production environment, demanding explicit user confirmation.

## The Conflict with `gr`

*Note: Some shell plugins (like Oh-My-Zsh git plugins) heavily alias `gr` to `git remote`. To avoid fatal git repository errors, the CLI binary is explicitly named `guardrail`.*

---

## 🚀 Installation

Ensure you have Go installed (1.24+). To install the CLI globally to your `$GOPATH/bin`:

```bash
make install
# or
go install github.com/ifeoluwashola/guardrail/...
```

> **Troubleshooting: `command not found`**  
> If you get `command not found: guardrail`, it means Go's bin directory isn't in your system's `PATH`. Run this command to add it to your profile:
>
> ```bash
> echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc && source ~/.zshrc
> ```

Verify the installation:

```bash
guardrail --help
```

---

## ⚙️ Configuration (`~/.guardrail/config.yaml`)

Guardrail is entirely config-driven. On first run, it looks for `~/.guardrail/config.yaml` in your local home directory. This file dictates your cloud and cluster contexts.

**Sample `~/.guardrail/config.yaml`:**

```yaml
environments:
  dev:
    cloud_profile: "aws-dev"
    kubernetes_context: "arn:aws:eks:us-east-1:123456789012:cluster/dev-cluster"
    terraform_workspace: "dev"
    is_production: false
  prod:
    cloud_profile: "aws-prod"
    kubernetes_context: "arn:aws:eks:us-east-1:123456789012:cluster/prod-cluster"
    terraform_workspace: "prod"
    is_production: true
```

---

## 🛠️ List of Commands & Usage

### 1. Context Switching (`guardrail use`)

Changes Guardrail's internal active context and executes shell integrations depending on the environment.

**Usage:**

```bash
guardrail use dev
guardrail use prod
```

*Note: If `is_production: true` is set, `use` will flash a bright red terminal warning letting you know you are now targeting production infrastructure.*

### 2. The Command Interceptor (`guardrail run -- <command>`)

This is the core safety engine. It wraps any command and executes it within the currently active environment. If the active environment is marked as **production**, it scans the command for high-risk keywords (`delete`, `apply`, `destroy`, `uninstall`).

If a risk is detected, it freezes the pipeline and requires manual exact-string validation of the target cluster.

**Safe Execution Passes Through:**

```bash
guardrail run -- kubectl get pods
```

**Dangerous Execution is Trapped:**

```bash
guardrail run -- kubectl delete deployment api-server
# DANGER: You are targeting PRODUCTION.
# ✔ Type the exact cluster context name ('arn:aws:eks...:cluster/prod-cluster') to confirm: 
```

### 3. Shell Profile Exporting (`guardrail env`)

A Go child-process cannot natively modify the `export` tree of your parent bash/zsh shell. To propagate environment variables seamlessly, evaluate the `env` command directly in your shell.

**Usage:**

```bash
eval $(guardrail env)
```

**Pro-Tip:** Add an alias to your `~/.bashrc` or `~/.zshrc` to automatically wrap the context execution!

```bash
alias gr-use='guardrail use $1 && eval $(guardrail env)'

# Now you can just type:
# gr-use dev
```

### 4. Custom PS1 Prompt (`guardrail prompt`)

Want to know exactly what environment you are operating in at all times? Insert Guardrail natively into your shell prompt layout.

**Usage:**

```bash
guardrail prompt
# outputs: [PROD | prod-cluster] (colored red)
```

**Pro-Tip (Zsh/Bash integration):**
Add this to your `~/.zshrc` profile:

```bash
export PS1="\$(guardrail prompt) \u@\h:\w$ "
```

### 5. Interactive Configuration (`guardrail set-context`)

Guardrail provides a suite of interactive commands so you never have to manually write YAML configs. **Pro-Tip: We've added short aliases like `sc` or `e` to save keystrokes!**

* `guardrail create-config` (aliases: `create`, `c`, `cc`): Interactively scaffolds your very first `~/.guardrail/config.yaml` using terminal prompts.
* `guardrail set-context` (aliases: `set`, `s`, `sc`): Spawns a CLI GUI dropdown to select an environment, allowing you to interactively override the `cloud_provider`, `cloud_profile`, K8s context, and production tags!
* `guardrail edit-config` (aliases: `edit`, `e`, `ec`): Opens your `~/.guardrail/config.yaml` dynamically in your host's `$EDITOR` (e.g. `vim`, `nano`).

---

## 🧪 Development & Testing

Guardrail uses standard Go tables for testing the safety interceptors.

```bash
make build   # builds local binary
make test    # runs cmd/run_test.go
```

---

## 📄 License

This project is open-source and licensed under the [MIT License](LICENSE). Feel free to fork, modify, and distribute it!
