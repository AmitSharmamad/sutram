create the following tables

1. deployment_log -> to track the progression of a deployment (immutable), 
add metadata JSONB column to this
    a. id
    b. deployment_id
    c. status
    d. metadata JSONB
    e. action_at
2. deployment_position -> to store the deployment positions
3. deployment_metrics -> to store the deployment metrics

fix the sql, repos, services and apis as needed