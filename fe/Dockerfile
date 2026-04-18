ARG NODE_VERSION=24.15.0
FROM node:${NODE_VERSION}-alpine
WORKDIR /app
COPY .output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
