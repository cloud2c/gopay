package alipay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloud2c/gopay"
)

// 商品图片上传接口 alipay.merchant.item.file.upload
// StatusCode = 200 is success
func (a *ClientV3) MerchantItemFileUpload(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantItemFileUploadRsp, err error) {
	err = bm.CheckEmptyError("scene", "file_content")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	res, bs, err := a.doProdPostFile(ctx, bm, v3MerchantItemFileUpload, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantItemFileUploadRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}
