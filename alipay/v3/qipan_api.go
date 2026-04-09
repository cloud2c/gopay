package alipay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloud2c/gopay"
)

// 上传创建人群 alipay.merchant.qipan.crowd.create
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanCrowdCreate(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanCrowdCreateRsp, err error) {
	err = bm.CheckEmptyError("crowd_name", "external_crowd_code", "user_list")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanCrowdCreate, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanCrowdCreateRsp{StatusCode: res.StatusCode}
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

// 人群中追加用户 alipay.merchant.qipan.crowduser.add
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanCrowdUserAdd(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanCrowdUserAddRsp, err error) {
	err = bm.CheckEmptyError("crowd_code", "user_list")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanCrowdUserAdd, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanCrowdUserAddRsp{StatusCode: res.StatusCode}
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

// 人群中删除用户 alipay.merchant.qipan.crowduser.delete
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanCrowdUserDelete(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanCrowdUserDeleteRsp, err error) {
	err = bm.CheckEmptyError("crowd_code", "user_list")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanCrowdUserDelete, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanCrowdUserDeleteRsp{StatusCode: res.StatusCode}
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

// 棋盘人群圈选标签基本信息查询 alipay.marketing.qipan.tagbase.batchquery
// StatusCode = 200 is success
func (a *ClientV3) MarketingQipanTagBaseBatchQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingQipanTagBaseBatchQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doGet(ctx, v3MarketingQipanTagBaseBatchQuery, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingQipanTagBaseBatchQueryRsp{StatusCode: res.StatusCode}
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

// 棋盘标签圈选值查询 alipay.marketing.qipan.tag.query
// StatusCode = 200 is success
func (a *ClientV3) MarketingQipanTagQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingQipanTagQueryRsp, err error) {
	err = bm.CheckEmptyError("tag_code")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3MarketingQipanTagQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingQipanTagQueryRsp{StatusCode: res.StatusCode}
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

// 查询人群列表 alipay.merchant.qipan.crowd.batchquery
// StatusCode = 200 is success
func (a *ClientV3) MarketingQipanCrowdBatchQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingQipanCrowdBatchQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingQipanCrowdBatchQuery, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingQipanCrowdBatchQueryRsp{StatusCode: res.StatusCode}
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

// 查询人群详情 alipay.merchant.qipan.crowd.query
// StatusCode = 200 is success
func (a *ClientV3) MarketingQipanCrowdQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingQipanCrowdQueryRsp, err error) {
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3MarketingQipanCrowdQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingQipanCrowdQueryRsp{StatusCode: res.StatusCode}
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

// 修改人群 alipay.merchant.qipan.crowd.modify
// StatusCode = 200 is success
func (a *ClientV3) MarketingQipanCrowdModify(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingQipanCrowdModifyRsp, err error) {
	err = bm.CheckEmptyError("crowd_code")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingQipanCrowdModify, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingQipanCrowdModifyRsp{StatusCode: res.StatusCode}
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

// 人群扩展接口 alipay.merchant.qipan.crowd.spread
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanCrowdSpread(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanCrowdSpreadRsp, err error) {
	err = bm.CheckEmptyError("apply_channel_list", "spread_count", "crowd_name", "seed_crowd_code", "is_include_seed")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanCrowdSpread, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanCrowdSpreadRsp{StatusCode: res.StatusCode}
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

// 上传创建灰黑产人群 alipay.merchant.qipan.greyblackcrowd.create
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanGreyBlackCrowdCreate(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanGreyBlackCrowdCreateRsp, err error) {
	err = bm.CheckEmptyError("crowd_name", "user_list")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanGreyBlackCrowdCreate, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanGreyBlackCrowdCreateRsp{StatusCode: res.StatusCode}
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

// 灰黑产人群中追加用户 alipay.merchant.qipan.greyblackcrowduser.add
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanGreyBlackCrowdUserAdd(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanGreyBlackCrowdUserAddRsp, err error) {
	err = bm.CheckEmptyError("crowd_code", "user_list")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanGreyBlackCrowdUserAdd, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanGreyBlackCrowdUserAddRsp{StatusCode: res.StatusCode}
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

// 灰黑产人群中删除用户 alipay.merchant.qipan.greyblackcrowduser.delete
// StatusCode = 200 is success
func (a *ClientV3) MerchantQipanGreyBlackCrowdUserDelete(ctx context.Context, bm gopay.BodyMap) (aliRsp *MerchantQipanGreyBlackCrowdUserDeleteRsp, err error) {
	err = bm.CheckEmptyError("crowd_code", "user_list")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MerchantQipanGreyBlackCrowdUserDelete, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MerchantQipanGreyBlackCrowdUserDeleteRsp{StatusCode: res.StatusCode}
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
