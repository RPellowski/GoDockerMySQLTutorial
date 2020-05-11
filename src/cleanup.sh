# remove containers
for container in frodo brandywine; do
  for cmd in kill rm; do
    docker ${cmd} ${container}
  done
done

# remove images
for image in hobbit shire; do
  for cmd in rmi; do
    docker ${cmd} ${image}
  done
done

# remove network
docker network rm my-net

# remove database persistent directory
sudo rm -fr ~/my-db

