// ACCESS_TOKEN=yourAccessToken,omitempty HTTP_PORT=8080 ENDPOINT=v1/tracking/user-activities k6 run --vus 100 --iterations 100000 stress-test/user-activities-post.js

import http from 'k6/http';

import { check } from 'k6';

const authorization = __ENV.ACCESS_TOKEN ? `Bearer ${__ENV.ACCESS_TOKEN}` : '';

export default function () {
  let res = http.post(
    `http://localhost:${__ENV.HTTP_PORT}/${__ENV.ENDPOINT}`,
    JSON.stringify({
      eventName: 'watch-15min',
      data: {
        itemId: "1",
      },
    }),
    {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': authorization,
      },
    },
  );

  check(res, { 'status was 200': (r) => r.status == 200 });
}
