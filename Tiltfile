k8s_yaml('infra/k8s/postgres.yaml')
k8s_yaml('infra/k8s/redis.yaml')
k8s_yaml('infra/k8s/auth-service.yaml')

docker_build(
    'auth-service',
    context='apps/auth-service',
    dockerfile='apps/auth-service/Dockerfile',
    live_update=[
        # Sync the local src folder to /app/src inside the container
        sync('apps/auth-service/src', '/app/src')
    ]
)

k8s_resource('auth-service', port_forwards='3001:3001')
