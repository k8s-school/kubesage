# KubeSage

**KubeSage** is a smart CLI assistant for DevOps and Kubernetes practitioners. It uses a Large Language Model (LLM) to:
- Audit and suggest improvements to Kubernetes YAML files
- Explain line-by-line what a YAML file does
- Generate new Kubernetes manifests from natural language prompts

## 📦 Features

- `audit`: Analyze and improve YAML files
- `--explain`: Output an explanation instead of an audit
- `--output`: Save the improved YAML to a file
- `generate`: Create YAML from natural language

## 🚀 Usage

### Audit a file
```bash
kubesage audit -f ./deployment.yaml -k $OPENAI_API_KEY
```

### Audit and save improved YAML
```bash
kubesage audit -f ./deployment.yaml -k $OPENAI_API_KEY -o improved.yaml
```

### Explain a file
```bash
kubesage audit -f ./deployment.yaml -k $OPENAI_API_KEY --explain
```

### Generate YAML from a prompt
```bash
kubesage generate -k $OPENAI_API_KEY "A deployment with 2 replicas of nginx exposing port 80"
```

## 🛠 Build

```bash
go mod init kubesage
go get github.com/spf13/cobra
go build -o kubesage
```

## 🧠 Powered by GPT-4 (or GPT-3.5)
You can use either model via OpenAI API.

## 📄 License

MIT
