import http from 'k6/http';
import { check } from 'k6';

export let options = {
  scenarios: {
    constant_rps: {
      executor: 'constant-arrival-rate',
      rate: 1000, // 1000 запросов в секунду
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 1500,
      maxVUs: 2000,
    },
  },
};

export default function () {
  let url = 'http://localhost:8080/api/v1/wallets/operation';
  let payload = JSON.stringify({
    walletID: '1d86baf8-d99e-4db4-9c2f-bb9b8ab4e2e1',
    operationType: 'DEPOSIT',
    amount: 1000,
  });

  let params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post(url, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
}
