module.exports = async ({ github, context }) => {
  const {
    TARGET_REPOSITORY,
    SERVICE_NAME,
  } = process.env;

  console.debug(JSON.stringify(context));

  await github.rest.repos.createDispatchEvent({
    owner: context.repo.owner,
    repo: TARGET_REPOSITORY,
    event_type: 'release-service',
    client_payload: { ref: context.ref, service: SERVICE_NAME },
  });
}
