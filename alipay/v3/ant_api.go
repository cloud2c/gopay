package alipay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloud2c/gopay"
)

// 蚂蚁店铺创建 ant.merchant.expand.shop.create
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantShopCreate(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantShopCreateRsp, err error) {
	err = bm.CheckEmptyError("business_address", "shop_category", "shop_type", "ip_role_id", "shop_name")
	if err != nil {
		return nil, err
	}
	if bm.GetString("contact_phone") == gopay.NULL && bm.GetString("contact_mobile") == gopay.NULL {
		return nil, errors.New("contact_phone and contact_mobile are not allowed to be null at the same time")
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3AntMerchantShopCreate, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantShopCreateRsp{StatusCode: res.StatusCode}
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

// 店铺查询接口 ant.merchant.expand.shop.query
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantShopQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantShopQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3AntMerchantShopQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantShopQueryRsp{StatusCode: res.StatusCode}
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

// 修改蚂蚁店铺 ant.merchant.expand.shop.modify
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantShopModify(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantShopModifyRsp, err error) {
	if bm.GetString("contact_phone") == gopay.NULL && bm.GetString("contact_mobile") == gopay.NULL {
		return nil, errors.New("contact_phone and contact_mobile are not allowed to be null at the same time")
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPatch(ctx, bm, v3AntMerchantShopModify, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantShopModifyRsp{StatusCode: res.StatusCode}
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

// 蚂蚁店铺关闭 ant.merchant.expand.shop.close
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantShopClose(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantShopCloseRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPatch(ctx, bm, v3AntMerchantShopClose, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantShopCloseRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 商户申请单查询 ant.merchant.expand.order.query
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantOrderQuery(ctx context.Context, orderId string, bm gopay.BodyMap) (aliRsp *AntMerchantOrderQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := fmt.Sprintf(v3AntMerchantOrderQuery, orderId) + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantOrderQueryRsp{StatusCode: res.StatusCode}
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

// 店铺分页查询接口 ant.merchant.expand.shop.page.query
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantShopPageQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantShopPageQueryRsp, err error) {
	err = bm.CheckEmptyError("ip_role_id", "page_num", "page_size")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3AntMerchantShopPageQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantShopPageQueryRsp{StatusCode: res.StatusCode}
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

// 图片上传 ant.merchant.expand.indirect.image.upload
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantExpandIndirectImageUpload(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantExpandIndirectImageUploadRsp, err error) {
	err = bm.CheckEmptyError("image_type", "image_content")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	res, bs, err := a.doProdPostFile(ctx, bm, v3AntMerchantExpandIndirectImageUpload, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantExpandIndirectImageUploadRsp{StatusCode: res.StatusCode}
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

// 商户mcc信息查询 ant.merchant.expand.mcc.query
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantExpandMccQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantExpandMccQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3AntMerchantExpandMccQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantExpandMccQueryRsp{StatusCode: res.StatusCode}
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

// 店铺增加收单账号 ant.merchant.expand.shop.receiptaccount.save
// StatusCode = 200 is success
func (a *ClientV3) AntMerchantExpandShopReceiptAccountSave(ctx context.Context, bm gopay.BodyMap) (aliRsp *AntMerchantExpandShopReceiptAccountSaveRsp, err error) {
	err = bm.CheckEmptyError("shop_id", "receipt_account_id")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3AntMerchantExpandShopReceiptAccountSave, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &AntMerchantExpandShopReceiptAccountSaveRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}
