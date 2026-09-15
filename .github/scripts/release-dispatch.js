module.exports = async ({ github, context }) => {
  const {
    TARGET_REPOSITORY,
    SERVICE_NAME,
  } = process.env;

  const release = context.payload.release;

  await github.rest.repos.createDispatchEvent({
    owner: context.repo.owner,
    repo: TARGET_REPOSITORY,
    event_type: 'image-published',
    client_payload: {
      service: SERVICE_NAME,
      tag: release.tag_name,
      release_url: release.html_url,
    },
  });
}
