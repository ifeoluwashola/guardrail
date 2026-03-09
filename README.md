# Guardrail 🛡️

Guardrail is a context-aware CLI tool written in Go that safely manages environment switching across various clusters and cloud accounts. Its core feature is a **Safety Interceptor** that halts execution when running high-risk commands (like `delete`, `uninstall`, or `destroy`) against a production environment, demanding explicit user confirmation.

**🌐 Full Documentation & Usage Guide:** [https://guardrail-website.vercel.app/](https://guardrail-website.vercel.app/)

## The Conflict with `gr`

*Note: Some shell plugins (like Oh-My-Zsh git plugins) heavily alias `gr` to `git remote`. To avoid fatal git repository errors, the CLI binary is explicitly named `guardrail`. However, you can alias it locally.*

---

## 🚀 Installation

Guardrail natively supports seamless distribution via GoReleaser across all major package managers.

### macOS / Linux (Homebrew)

```bash
brew tap ifeoluwashola/homebrew-tap
brew install guardrail
```

### Debian / Ubuntu (APT)

```bash
curl -LO https://github.com/ifeoluwashola/guardrail/releases/latest/download/guardrail_Linux_x86_64.deb
sudo dpkg -i guardrail_Linux_x86_64.deb
```

### Windows (Scoop)

```powershell
scoop bucket add guardrail https://github.com/ifeoluwashola/scoop-bucket.git
scoop install guardrail
```

### Compile From Source (Requires Go 1.24+)

```bash
go install github.com/ifeoluwashola/guardrail@latest
```

---

## ⚙️ Configuration (`~/.guardrail/config.yaml`)

Guardrail is entirely config-driven. On first run, it looks for `~/.guardrail/config.yaml`. Provide definitions for your cloud and cluster contexts.

```yaml
environments:
  dev:
    kubernetes_context: "arn:aws:eks:us-east-1:123456789012:cluster/dev-cluster"
    is_production: false
  prod:
    kubernetes_context: "arn:aws:eks:us-east-1:123456789012:cluster/prod-cluster"
    is_production: true
```

---

## 🛠️ List of Commands & Usage

Guardrail features interactive configuration wizards, meaning you never actually have to touch the YAML file manually!

### 1. Initialization & Configuration

* **`guardrail create-config`** (alias: `cc`): Scaffolds the initial base config.
* **`guardrail set-context`** (alias: `sc`): Interactive terminal dropdown UI to add or modify environments.
* **`guardrail edit-config`**: Falls back to your local `$EDITOR` for manual YAML tuning.

### 2. Environment Management

* **`guardrail list-contexts`** (alias: `lc`): Prints all environments configured in your YAML, brilliantly highlighting your *Active* selection and explicitly flagging any [PROD] environments.
* **`guardrail use [env_name]`**: Switches the internal state to the requested cluster target.

### 3. Execution & Safety Interceptor

The core safety engine. Wrap your command with `guardrail run` to execute beneath the active profile.

```bash
guardrail run -- kubectl delete pods --all
```

If your active context has `is_production: true`, it physically halts execution for any destructive keyword (`delete`, `apply`, `destroy`) and demands the exact cluster string to proceed. **Zero accidental production drops.**

### 4. Shell Profile Interfacing

* **`guardrail env`**: Prints out the bash `export` statements required to inject the context into your local shell.
* **`guardrail prompt`**: Renders a highly-visible badge (e.g. `[PROD | my-cluster]`), built perfectly for custom Zsh or Starship themes!

**👉 For deep dives into Starship (`ps1`) integration instructions or Shell Aliasing, please visit the [Official Documentation](https://guardrail-website.vercel.app/).**

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
