name: Publish Docker image
on:
  push:
    branches:
      - production
jobs:
  push_to_registry:
    name: Push Docker image to GitHub Packages
    runs-on: ubuntu-latest
    steps:
      - name: Check out the repo
        uses: actions/checkout@v2
      - name: Push to GitHub Packages
        uses: docker/build-push-action@v1
        with:
          username: ${{ github.actor }}
          password: ${{ secrets.REGISTRY_TOKEN }}
          registry: ghcr.io
          repository: sindreslungaard/duel-masters/production
          tags: latest
