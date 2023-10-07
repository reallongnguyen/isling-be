// HTTP_PORT=8080 ENDPOINT=stress-test k6 run --vus 100 --iterations 100000 stress-test/empty-post.js

import http from 'k6/http';

import { check } from 'k6';

export default function () {
  let res = http.post(
    `http://localhost:${__ENV.HTTP_PORT}/${__ENV.ENDPOINT}`,
    JSON.stringify({
      title: 'My awesome test',
      description: 'This is a test',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );

  check(res, { 'status was 200': (r) => r.status == 200 });
}
