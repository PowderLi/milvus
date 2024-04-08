// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpcclient

import (
	"bytes"
	"context"
	"github.com/milvus-io/milvus/pkg/util/compressor/datadog"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/milvus-io/milvus-proto/go-api/v2/milvuspb"
	"github.com/milvus-io/milvus/pkg/util/compressor/deflate"
	"github.com/milvus-io/milvus/pkg/util/compressor/lz4"
	"github.com/milvus-io/milvus/pkg/util/compressor/snappy"
	"github.com/milvus-io/milvus/pkg/util/compressor/zstd"
	"github.com/milvus-io/milvus/pkg/util/paramtable"
)

var compressorTypes = []string{zstd.Name, gzip.Name, snappy.Name, lz4.Name, deflate.Name}

func TestGrpcEncoder(t *testing.T) {
	data := "hello zstd algorithm!"

	for _, name := range compressorTypes {
		compressor := encoding.GetCompressor(name)
		var buf bytes.Buffer
		writer, err := compressor.Compress(&buf)
		assert.NoError(t, err)
		written, err := writer.Write([]byte(data))
		assert.NoError(t, err)
		assert.Equal(t, written, len(data))
		err = writer.Close()
		assert.NoError(t, err)

		reader, err := compressor.Decompress(bytes.NewReader(buf.Bytes()))
		assert.NoError(t, err)
		result := make([]byte, len(data))
		reader.Read(result)
		assert.Equal(t, data, string(result))
	}
}

const (
	data   = "// Licensed to the LF AI & Data foundation under one\n// or more contributor license agreements. See the NOTICE file\n// distributed with this work for additional information\n// regarding copyright ownership. The ASF licenses this file\n// to you under the Apache License, Version 2.0 (the\n// \"License\"); you may not use this file except in compliance\n// with the License. You may obtain a copy of the License at\n//\n//\thttp://www.apache.org/licenses/LICENSE-2.0\n//\n// Unless required by applicable law or agreed to in writing, software\n// distributed under the License is distributed on an \"AS IS\" BASIS,\n// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n// See the License for the specific language governing permissions and\n// limitations under the License."
	vector = "[0.0459, 0.0439, 0.0251, 0.0318, 0.0159, 0.0564, 0.0335, 0.0273, 0.0156, 0.0627, 0.0504, 0.0116, 0.0289, 0.0432, 0.0388, 0.0282, 0.0405, 0.0417, 0.0309, 0.0338, 0.0165, 0.0291, 0.0245, 0.0208, 0.0207, 0.0727, 0.0386, 0.0145, 0.0347, 0.0462, 0.0238, 0.0333, 0.0616, 0.0418, 0.0344, 0.0448, 0.0221, 0.0348, 0.0275, 0.0319, 0.0445, 0.1036, 0.0365, 0.0168, 0.0539, 0.0554, 0.0224, 0.0432, 0.1612, 0.0764, 0.0892, 0.1059, 0.0974, 0.057, 0.064, 0.072, 0.1108, 0.1132, 0.0399, 0.035, 0.0914, 0.0654, 0.0676, 0.0481, 0.2254, 0.0976, 0.1954, 0.1424, 0.1787, 0.0932, 0.0989, 0.0909, 0.1794, 0.0632, 0.104, 0.0889, 0.1972, 0.0874, 0.1693, 0.1058, 0.0982, 0.0606, 0.1163, 0.0678, 0.0728, 0.065, 0.029, 0.0384, 0.0692, 0.0409, 0.0631, 0.0669, 0.1183, 0.0573, 0.0792, 0.1229, 0.0317, 0.0608, 0.0381, 0.0357, 0.0282, 0.0554, 0.0146, 0.0199, 0.0332, 0.0343, 0.0357, 0.0265, 0.0308, 0.0283, 0.0285, 0.0467, 0.0312, 0.0634, 0.0273, 0.0298, 0.0189, 0.0766, 0.0253, 0.0241, 0.0218, 0.0432, 0.0398, 0.0132, 0.0252, 0.0278, 0.0262, 0.0261, 0.0442, 0.0496, 0.0578, 0.0428, 0.0215, 0.077, 0.0415, 0.0274, 0.0247, 0.0802, 0.0313, 0.0157, 0.0308, 0.0564, 0.0409, 0.0171, 0.0445, 0.0431, 0.0398, 0.0338, 0.015, 0.0268, 0.0264, 0.0234, 0.0342, 0.0785, 0.044, 0.0111, 0.0373, 0.0732, 0.0441, 0.0267, 0.062, 0.0705, 0.052, 0.0588, 0.0347, 0.0398, 0.0434, 0.0626, 0.0348, 0.0771, 0.0439, 0.0277, 0.0688, 0.1063, 0.0426, 0.0315, 0.1182, 0.108, 0.1233, 0.1331, 0.1333, 0.0894, 0.084, 0.147, 0.1676, 0.0873, 0.0449, 0.0598, 0.1058, 0.1296, 0.1083, 0.0598, 0.3077, 0.171, 0.2578, 0.1343, 0.3025, 0.148, 0.1543, 0.115, 0.3003, 0.129, 0.1667, 0.1124, 0.2835, 0.1885, 0.268, 0.1088, 0.1394, 0.0846, 0.1385, 0.0723, 0.1257, 0.1025, 0.0347, 0.0402, 0.096, 0.0779, 0.0715, 0.1003, 0.1399, 0.1198, 0.1459, 0.1709, 0.0286, 0.0901, 0.0579, 0.0447, 0.0259, 0.072, 0.027, 0.0172, 0.0401, 0.0367, 0.0308, 0.0446, 0.0333, 0.0428, 0.0485, 0.0842, 0.0423, 0.0993, 0.0507, 0.0226, 0.0334, 0.1118, 0.0455, 0.0215, 0.034, 0.0626, 0.0287, 0.0151, 0.0331, 0.0358, 0.0273, 0.0148, 0.0481, 0.0572, 0.1005, 0.041, 0.0327, 0.039, 0.0301, 0.0316, 0.0348, 0.0455, 0.0162, 0.0291, 0.0381, 0.0462, 0.0711, 0.0388, 0.1661, 0.0992, 0.131, 0.1682, 0.2028, 0.1029, 0.1002, 0.1855, 0.1282, 0.0897, 0.0436, 0.0749, 0.1315, 0.1428, 0.0925, 0.0733, 0.3956, 0.2322, 0.2294, 0.3246, 0.4046, 0.2421, 0.1451, 0.292, 0.3793, 0.2134, 0.1256, 0.2624, 0.3962, 0.2613, 0.2024, 0.2949, 0.1559, 0.1304, 0.1246, 0.1057, 0.1153, 0.1169, 0.0454, 0.0656, 0.1537, 0.0689, 0.0724, 0.1378, 0.1997, 0.1151, 0.1348, 0.1879, 0.0378, 0.0295, 0.0267, 0.0348, 0.0299, 0.1118, 0.0764, 0.0252, 0.0216, 0.0718, 0.0883, 0.024, 0.0255, 0.0446, 0.0302, 0.0252, 0.0384, 0.0321, 0.0286, 0.0506, 0.0365, 0.0743, 0.0723, 0.0205, 0.0279, 0.1004, 0.076, 0.0191, 0.0269, 0.0528, 0.0361, 0.0294, 0.0602, 0.0421, 0.0368, 0.0597, 0.0258, 0.0539, 0.0689, 0.023, 0.05, 0.1801, 0.0726, 0.0224, 0.0318, 0.051, 0.028, 0.0313, 0.1244, 0.0621, 0.0508, 0.0942, 0.077, 0.0747, 0.1127, 0.0799, 0.0912, 0.1698, 0.0911, 0.05, 0.1001, 0.049, 0.0347, 0.0251, 0.1328, 0.0812, 0.1235, 0.0955, 0.1485, 0.0946, 0.1559, 0.0929, 0.1354, 0.0936, 0.1809, 0.1137, 0.1718, 0.0684, 0.12, 0.0801, 0.0574, 0.0528, 0.0601, 0.0403, 0.0676, 0.0763, 0.0613, 0.0456, 0.0588, 0.0613, 0.1312, 0.0859, 0.0948, 0.0499, 0.052, 0.0911, 0.0325, 0.0321, 0.025, 0.0217, 0.0272, 0.068, 0.0372, 0.0287, 0.0296, 0.0501, 0.0738, 0.0348, 0.0293, 0.0277, 0.0336, 0.0375, 0.0297, 0.0363, 0.0249, 0.0299, 0.0266, 0.0953, 0.0633, 0.0312, 0.0307, 0.0476, 0.0782, 0.0227, 0.024, 0.0336, 0.0265, 0.023, 0.035, 0.0365, 0.0427, 0.0366, 0.0246, 0.0864, 0.0795, 0.0218, 0.0382, 0.0744, 0.0657, 0.0357, 0.0228, 0.0661, 0.0471, 0.0146, 0.0397, 0.0349, 0.0395, 0.0464, 0.0203, 0.0447, 0.0729, 0.0205, 0.0433, 0.075, 0.083, 0.0212, 0.04, 0.0859, 0.0581, 0.0338, 0.0453, 0.0414, 0.0623, 0.0866, 0.0264, 0.0572, 0.0933, 0.0392, 0.0314, 0.109, 0.092, 0.0258, 0.0393, 0.0864, 0.0511, 0.037, 0.0811, 0.074, 0.0903, 0.1156, 0.092, 0.1179, 0.1505, 0.1201, 0.1061, 0.1405, 0.1053, 0.0583, 0.0773, 0.0845, 0.0484, 0.026, 0.2203, 0.1299, 0.1846, 0.0925, 0.2102, 0.1497, 0.2352, 0.0999, 0.188, 0.1712, 0.2815, 0.1449, 0.1963, 0.1609, 0.1918, 0.0808, 0.0825, 0.0618, 0.0858, 0.0531, 0.0997, 0.0804, 0.0711, 0.0411, 0.0658, 0.1118, 0.1539, 0.1242, 0.096, 0.1132, 0.0893, 0.122, 0.0253, 0.0549, 0.0244, 0.0252, 0.0308, 0.0559, 0.0656, 0.0236, 0.0404, 0.0322, 0.0613, 0.0598, 0.0347, 0.0375, 0.039, 0.0451, 0.0408, 0.0565, 0.0381, 0.0251, 0.033, 0.1151, 0.089, 0.0261, 0.0415, 0.0602, 0.0762, 0.0268, 0.032, 0.0312, 0.0475, 0.0221, 0.0496, 0.0325, 0.0585, 0.0332, 0.0408, 0.0616, 0.0698, 0.0395, 0.0429, 0.0651, 0.082, 0.0386, 0.0405, 0.0342, 0.0579, 0.0271, 0.0939, 0.0647, 0.0884, 0.0954, 0.14, 0.1124, 0.1252, 0.105, 0.1049, 0.1092, 0.0498, 0.0511, 0.1082, 0.1031, 0.0545, 0.0474, 0.2339, 0.1581, 0.1452, 0.1777, 0.2803, 0.2541, 0.1426, 0.1622, 0.2905, 0.2539, 0.141, 0.1648, 0.2805, 0.2083, 0.1199, 0.1783, 0.0945, 0.072, 0.0748, 0.0624, 0.077, 0.116, 0.0494, 0.0385, 0.1258, 0.0837, 0.095, 0.0853, 0.1383, 0.0923, 0.0896, 0.1217, 0.0219, 0.0325, 0.0252, 0.0218, 0.0134, 0.0766, 0.057, 0.0276, 0.024, 0.0582, 0.0787, 0.0289, 0.021, 0.0314, 0.028, 0.0213, 0.0291, 0.0395, 0.0382, 0.0296, 0.0175, 0.0512, 0.0698, 0.0267, 0.0245, 0.0752, 0.0793, 0.0232, 0.0263, 0.0335, 0.0342, 0.0296, 0.0329, 0.0374, 0.029, 0.0595, 0.0235, 0.043, 0.0518, 0.0456, 0.0474, 0.111, 0.0492, 0.0248, 0.0317, 0.0458, 0.0297, 0.0271, 0.0739, 0.0451, 0.0458, 0.0813, 0.0789, 0.0538, 0.0909, 0.0962, 0.0869, 0.1304, 0.0801, 0.0452, 0.0858, 0.0506, 0.0496, 0.0315, 0.0812, 0.0385, 0.0875, 0.0608, 0.131, 0.0597, 0.1399, 0.0882, 0.1321, 0.0864, 0.1482, 0.0969, 0.1444, 0.0704, 0.0887, 0.0552, 0.0418, 0.0448, 0.0561, 0.0491, 0.0658, 0.0617, 0.057, 0.0653, 0.0605, 0.0585, 0.0923, 0.0708, 0.0723, 0.0571, 0.0414, 0.0629, 0.0273, 0.0428, 0.0371, 0.0306, 0.0345, 0.0721, 0.0302, 0.0345, 0.0224, 0.052, 0.0572, 0.0428, 0.0288, 0.0317, 0.0263, 0.0414, 0.0285, 0.0348, 0.0291, 0.0282, 0.0206, 0.0705, 0.049, 0.0296, 0.0214, 0.0694, 0.0656, 0.0416, 0.0306, 0.0256, 0.0202, 0.0232]"
)

func benchmarkRun(b *testing.B, compressionName, data string) {
	compressor := encoding.GetCompressor(compressionName)

	// Reset the timer to exclude setup time from the measurements
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			var buf bytes.Buffer
			writer, _ := compressor.Compress(&buf)
			writer.Write([]byte(data))
			writer.Close()

			compressedData := buf.Bytes()
			reader, _ := compressor.Decompress(bytes.NewReader(compressedData))
			var result bytes.Buffer
			result.ReadFrom(reader)
		}
	}
}

func BenchmarkVarcharDeflate(b *testing.B) {
	benchmarkRun(b, deflate.Name, data)
}
func BenchmarkVarcharLz4(b *testing.B) {
	benchmarkRun(b, lz4.Name, data)
}
func BenchmarkVarcharSnappy(b *testing.B) {
	benchmarkRun(b, snappy.Name, data)
}
func BenchmarkVarcharZstd(b *testing.B) {
	benchmarkRun(b, zstd.Name, data)
}
func BenchmarkVarcharDatadog(b *testing.B) {
	benchmarkRun(b, datadog.Name, data)
}
func BenchmarkVarcharGzip(b *testing.B) {
	benchmarkRun(b, gzip.Name, data)
}

func BenchmarkVectorDeflate(b *testing.B) {
	benchmarkRun(b, deflate.Name, vector)
}
func BenchmarkVectorLz4(b *testing.B) {
	benchmarkRun(b, lz4.Name, vector)
}
func BenchmarkVectorSnappy(b *testing.B) {
	benchmarkRun(b, snappy.Name, vector)
}
func BenchmarkVectorZstd(b *testing.B) {
	benchmarkRun(b, zstd.Name, vector)
}
func BenchmarkVectorDatadog(b *testing.B) {
	benchmarkRun(b, datadog.Name, vector)
}
func BenchmarkVectorGzip(b *testing.B) {
	benchmarkRun(b, gzip.Name, vector)
}

func TestGrpcCompression(t *testing.T) {
	// ClusterInjectionUnaryClientInterceptor/ClusterInjectionStreamClientInterceptor need read `msgChannel.chanNamePrefix.cluster`
	paramtable.Init()
	lis, err := net.Listen("tcp", "localhost:")
	address := lis.Addr()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &compressionServer{})
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	for _, compressorType := range compressorTypes {
		_, err := sayHelloWithCompression(address.String(), compressorType)
		assert.NoError(t, err)
	}

	invalidCompression := "invalid"
	_, err = sayHelloWithCompression(address.String(), invalidCompression)
	assert.Equal(t, true, strings.HasSuffix(err.Error(),
		status.Errorf(codes.Internal, "grpc: Compressor is not installed for requested grpc-encoding \"%s\"", invalidCompression).Error()))
}

func sayHelloWithCompression(serverAddress string, compressionType string) (any, error) {
	clientBase := ClientBase[grpcClientInterface]{
		ClientMaxRecvSize:      1 * 1024 * 1024,
		ClientMaxSendSize:      1 * 1024 * 1024,
		DialTimeout:            60 * time.Second,
		KeepAliveTime:          60 * time.Second,
		KeepAliveTimeout:       60 * time.Second,
		RetryServiceNameConfig: "helloworld.Greeter",
		MaxAttempts:            1,
		InitialBackoff:         10.0,
		MaxBackoff:             60.0,
		CompressionEnabled:     true,
		CompressionType:        compressionType,
	}
	clientBase.SetGetAddrFunc(func() (string, error) {
		return serverAddress, nil
	})
	clientBase.SetNewGrpcClientFunc(func(cc *grpc.ClientConn) grpcClientInterface {
		return &compressionClient{helloworld.NewGreeterClient(cc)}
	})
	defer clientBase.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return clientBase.Call(ctx, func(client grpcClientInterface) (any, error) {
		return client.(*compressionClient).SayHello(ctx, &helloworld.HelloRequest{})
	})
}

type grpcClientInterface interface {
	GetComponentStates(ctx context.Context, in *milvuspb.GetComponentStatesRequest, opts ...grpc.CallOption) (*milvuspb.ComponentStates, error)
}

type compressionClient struct {
	helloworld.GreeterClient
}

func (c *compressionClient) GetComponentStates(ctx context.Context, in *milvuspb.GetComponentStatesRequest, opts ...grpc.CallOption) (*milvuspb.ComponentStates, error) {
	return &milvuspb.ComponentStates{}, nil
}

type compressionServer struct {
	helloworld.UnimplementedGreeterServer
}

func (compressionServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}
