package main

// func TestHandleWalletBalance(t *testing.T) {
// 	mockClient := mocks.NewLightningClient(t)

// 	mockClient.On("WalletBalance", mock.Anything, mock.Anything).
// 		Return(&lnrpc.WalletBalanceResponse{TotalBalance: 21}, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	res := httptest.NewRecorder()

// 	handleWalletBalance(res, req, mockClient)

// 	bodyBytes, err := io.ReadAll(res.Result().Body)
// 	assert.Nil(t, err)
// 	assert.Contains(t, string(bodyBytes), "21")
// }

// func TestHandleNewAddress(t *testing.T) {
// 	mockClient := mocks.NewLightningClient(t)

// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	res := httptest.NewRecorder()

// 	addressStr := "bcrt1qddzehdyj5e7w4sfya3h9qznnm80etc9gkpk0qd"
// 	mockClient.On("NewAddress", mock.Anything, mock.Anything).
// 		Return(&lnrpc.NewAddressResponse{Address: addressStr}, nil)

// 	handleNewAddress(res, req, mockClient)

// 	bodyBytes, err := io.ReadAll(res.Result().Body)
// 	assert.Nil(t, err)
// 	assert.Contains(t, string(bodyBytes), addressStr)
// }

// func TestHandlePubKey(t *testing.T) {
// 	mockClient := mocks.NewLightningClient(t)

// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	res := httptest.NewRecorder()

// 	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
// 	mockClient.On("GetInfo", mock.Anything, mock.Anything).
// 		Return(&lnrpc.GetInfoResponse{IdentityPubkey: pubKey}, nil)

// 	handlePubKey(res, req, mockClient)

// 	bodyBytes, err := io.ReadAll(res.Result().Body)
// 	assert.Nil(t, err)
// 	assert.Contains(t, string(bodyBytes), pubKey)
// }

// func TestHandleConnectPeer(t *testing.T) {
// 	mockClient := mocks.NewLightningClient(t)
// 	assert := assert.New(t)

// 	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
// 	host := "lnd1.myfancysats.com"
// 	port := 9745

// 	connectPeerReq := types.ConnectPeerReq{PubKey: pubKey, Host: host, Port: port}
// 	connectPeerBytes, err := json.Marshal(connectPeerReq)
// 	assert.Nil(err)

// 	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(connectPeerBytes))
// 	res := httptest.NewRecorder()

// 	mockClient.On("ConnectPeer", mock.Anything, mock.Anything).
// 		Return(&lnrpc.ConnectPeerResponse{}, nil)

// 	handleConnectPeer(res, req, mockClient)

// 	bodyBytes, err := io.ReadAll(res.Result().Body)
// 	assert.Nil(err)
// 	assert.Contains(string(bodyBytes), "request received")
// }

// func TestHandleOpenChannel(t *testing.T) {
// 	mockClient := mocks.NewLightningClient(t)
// 	assert := assert.New(t)

// 	pubKey := "037c70cddec9b27c92af73a6b04cf09672fb29b18eca86890d835779979ff61c40"
// 	localAmt := 20001

// 	openChannelReq := types.OpenChannelReq{PubKey: pubKey, LocalAmtSats: uint64(localAmt)}
// 	openChannelReqBytes, err := json.Marshal(openChannelReq)
// 	assert.Nil(err)

// 	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(openChannelReqBytes))
// 	res := httptest.NewRecorder()

// 	mockClient.On("OpenChannelSync", mock.Anything, mock.Anything).
// 		Return(&lnrpc.ChannelPoint{OutputIndex: uint32(615)}, nil)

// 	handleOpenChannel(res, req, mockClient)

// 	bodyBytes, err := io.ReadAll(res.Result().Body)
// 	assert.Nil(err)
// 	assert.Contains(string(bodyBytes), "615")
// }
