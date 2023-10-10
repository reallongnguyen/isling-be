// ENDPOINT=stress-test k6 run --vus 100 --iterations 100000 stress-test/empty-get.js

import http from 'k6/http';

import { check } from 'k6';

const host = __ENV.HTTP_HOST || 'http://localhost';
const post = __ENV.HTTP_PORT || '8080';

export default function () {
  let res = http.get(
    `${host}:${post}/${__ENV.ENDPOINT}`,
    undefined,
    { headers: { 'Content-Type': 'application/json' } }
  );

  check(res, { 'status was 200': (r) => r.status == 200 });
}
