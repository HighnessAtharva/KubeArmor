package core

import (
	"context"
	"math/rand"
	"testing"

	pb "github.com/accuknox/KubeArmor/LogServer/protobuf"
	"google.golang.org/grpc"
)

func TestLogServer(t *testing.T) {
	Output = false
	serverAddr := "localhost:32767"

	// == //

	t.Log("[INFO] Start LogServer")

	server := NewLogServer(":32767")

	t.Log("[PASS] Started LogServer")

	// == //

	t.Log("[INFO] Start to receive logs")

	go server.ReceiveLogs()

	t.Log("[PASS] Started to receive logs")

	// == //

	t.Log("[INFO] Connect to LogServer")

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}
	t.Log("[PASS] Connected to LogServer")

	// == //

	t.Log("[INFO] Create a log message client")

	client := pb.NewLogMessageClient(conn)

	t.Log("[PASS] Created the log message client")

	// == //

	t.Log("[INFO] Test HealthCheck API")

	rand := rand.Int31()
	nonce := pb.NonceMessage{Nonce: rand}
	res, err := client.HealthCheck(context.Background(), &nonce)
	if err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}

	if rand != res.Retval {
		t.Errorf("[FAIL] Unexpected result: %d (expected: %d)", res.Retval, rand)
		return
	}
	t.Log("[PASS] Tested HealthCheck API")

	// == //

	t.Log("[INFO] Test AuditLogs API")

	auditStream, err := client.AuditLogs(context.Background())
	if err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}

	auditLog := pb.AuditLog{}
	auditStream.Send(&auditLog)

	res, err = auditStream.CloseAndRecv()
	if err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}

	if res.Retval != 0 {
		t.Errorf("[FAIL] Unexpected result: %d", res.Retval)
	}
	t.Log("[PASS] Tested AuditLogs API")

	// == //

	t.Log("[INFO] Test SystemLogs API")

	systemStream, err := client.SystemLogs(context.Background())
	if err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}

	systemLog := pb.SystemLog{}
	systemStream.Send(&systemLog)

	res, err = systemStream.CloseAndRecv()
	if err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}

	if res.Retval != 0 {
		t.Errorf("[FAIL] Unexpected result: %d", res.Retval)
	}
	t.Log("[PASS] Tested SystemLogs API")

	// == //

	t.Log("[INFO] Disconnect from LogServer")

	if err := conn.Close(); err != nil {
		t.Errorf("[FAIL] %v", err)
		return
	}
	t.Log("[PASS] Disconnected from LogServer")

	// == //

	t.Log("[INFO] Stop LogServer")

	server.DestroyLogServer()

	t.Log("[PASS] Stopped LogServer")

	// == //
}
