module.exports = async ({ github, context }) => {
  const {
    TARGET_REPOSITORY,
    SERVICE_NAMES,
  } = process.env;

  const release = context.payload.release;

  const services = JSON.parse(SERVICE_NAMES);

  await github.rest.repos.createDispatchEvent({
    owner: context.repo.owner,
    repo: TARGET_REPOSITORY,
    event_type: 'image-published',
    client_payload: {
      services,
      tag: release.tag_name,
      release_url: release.html_url,
    },
  });
}
