# KrakenD based LLM model router

This is a [KrakenD](https://www.krakend.io/docs/ai-gateway/) based LLM model router implementation. This is an egress API gateway meant to route between different LLM providers, such as OpenAI or Gemini.

### Building and Running

```bash
# install local development dependencies
brew bundle --no-lock --verbose

# Build the Docker image
docker build -t agentic-layer/model-router-krakend .

# Run the container with required environment variables
docker run -p 8080:8080 -e OPENAI_API_KEY=$OPENAI_API_KEY agentic-layer/model-router-krakend
```

## Deployment

```bash
# in order for the proxy to work we have to manually create a Kubernetes secrets
# that contains an OPENAI_API_KEY environment variable
# that contains an ANTHROPIC_API_KEY environment variable
kubectl create secret generic openai-api-key --from-literal=OPENAI_API_KEY=$OPENAI_API_KEY
kubectl apply -k kustomize/local/

# to test the proxy, issue the following curl command
curl http://krakend.127.0.0.1.sslip.io/llm/openai \
  -H "Content-Type: application/json" \
  -d '{
     "model": "gpt-4o-mini",
     "messages": [{"role": "user", "content": "Say this is a test!"}],
     "temperature": 0.7
   }'
```

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the Apache v2.0 open source license, read the `LICENSE` file for details.
