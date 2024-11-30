if [ "$USE_PROFILE" == "prod" ]; then
  cp /build/.env.prod /app/.env
  cp /build/firebase.json /app/
  echo "use prod"
elif [ "$USE_PROFILE" == "live" ]; then
  cp /build/.env.live /app/.env
  cp /build/test-firebase.json /app/
  echo "use live"
fi

exec /bin/breezenote