k8s_yaml('infra/k8s/postgres.yaml')
k8s_yaml('infra/k8s/redis.yaml')
k8s_yaml('infra/k8s/auth-service.yaml')
k8s_yaml('infra/k8s/api-gateway.yaml')

docker_build(
    'auth-service',
    context='apps/auth-service',
    dockerfile='apps/auth-service/Dockerfile',
    live_update=[
        # Sync the local src folder to /app/src inside the container
        sync('apps/auth-service/src', '/app/src')
    ]
)

docker_build(
    'api-gateway',
    context='apps/api-gateway',
    dockerfile='apps/api-gateway/Dockerfile',
    live_update=[
        # Sync the local src folder to /app/src inside the container
        sync('apps/api-gateway', '/app')
    ]
)

k8s_resource('auth-service', port_forwards='3001:3001')
k8s_resource('api-gateway')
k8s_resource('postgres-db', port_forwards='5434:5432')
k8s_resource('redis-cache', port_forwards='63799:6379')
