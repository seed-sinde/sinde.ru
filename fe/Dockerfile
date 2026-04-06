ARG BUN_VERSION=1.3.11
FROM oven/bun:${BUN_VERSION}-alpine
WORKDIR /app
COPY .output ./.output
ENV NODE_ENV=development
ENV HOST=0.0.0.0
ENV PORT=3000
EXPOSE 3000
CMD ["bun", ".output/server/index.mjs"]
