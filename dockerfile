FROM node:22-alpine AS build-env
COPY . /app
WORKDIR /app


















FROM gcr.io/distroless/nodejs22-debian13 AS runner
COPY --from=build-env /app /app
WORKDIR /app
