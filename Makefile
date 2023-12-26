all:
	go build -o bin/circles

css:
	pnpm css

dev:
	air & pnpm dev
