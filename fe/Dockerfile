ARG BUN_VERSION=1.3.12
FROM oven/bun:${BUN_VERSION}-alpine
WORKDIR /app
COPY .output ./.output
EXPOSE 3000
CMD ["bun", ".output/server/index.mjs"]
