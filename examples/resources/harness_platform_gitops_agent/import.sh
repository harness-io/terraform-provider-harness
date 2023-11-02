# Import an Account level Gitops Agent
terraform import harness_platform_gitops_agent.example <agent_id>

# Import an Org level Gitops Agent
terraform import harness_platform_gitops_agent.example <organization_id>/<agent_id>

# Import a Project level Gitops Agent
terraform import harness_platform_gitops_agent.example <organization_id>/<project_id>/<agent_id>
