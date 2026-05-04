import grpc from 'k6/net/grpc';
import { check } from 'k6';

const GRPC_ADDRESS = __ENV.GRPC_ADDRESS || "app:50051";
const DURATION = __ENV.DURATION || "1m";
const VUS = parseInt(__ENV.VUS) || 500;

const client = new grpc.Client();
client.load(['/proto'], 'gotips/v1/pcpart.proto');
client.load(['/proto'], 'gotips/v1/pcpart_store_service.proto');

export const options = {
    stages: [
        { duration: '30s', target: 10 },
        { duration: DURATION, target: VUS },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        'checks': ['rate>0.9'],
        'grpc_req_duration': ['p(95)<500'],
    },
};

// Переменная внутри VU (не глобальная для всех, а своя у каждого VU)
let connected = false;

export default () => {
    // Если этот конкретный VU еще не подключен — подключаемся
    if (!connected) {
        client.connect(GRPC_ADDRESS, {
            plaintext: true,
            timeout: '5s' // Даем время на установку соединения
        });
        connected = true;
    }

    /*
    const response = client.invoke('gotips.v1.PcPartStoreService/GetPcPart', {
        "id": "019d5da1-46dd-7b0d-82e5-49345ac87e79"
    });
     */
    const response = client.invoke('gotips.v1.PcPartStoreService/GetPcPartsRecent', {"limit": "LIMIT_SMALL"});

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    // МЫ НЕ ЗАКРЫВАЕМ СОЕДИНЕНИЕ (client.close() удален)
};