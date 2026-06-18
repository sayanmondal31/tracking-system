FROM oven/bun:1.0.5-alpine AS builder

WORKDIR /app

# copy dependencies
COPY package.json .

# install dependencies including devDependencies required for compilation)
RUN bun install 

FROM oven/bun:1.0.5-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

COPY --from=builder /app/node_modules ./node_modules

COPY . .

EXPOSE 3001

CMD ["bun","run","src/server.ts"]






