package alipay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloud2c/gopay"
)

// 会员卡模板创建 alipay.marketing.card.template.create
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardTemplateCreate(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardTemplateCreateRsp, err error) {
	err = bm.CheckEmptyError("request_id", "template_style_info")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingCardTemplateCreate, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardTemplateCreateRsp{StatusCode: res.StatusCode}
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

// 会员卡模板查询接口 alipay.marketing.card.template.query
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardTemplateQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardTemplateQueryRsp, err error) {
	err = bm.CheckEmptyError("template_id")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	uri := v3MarketingCardTemplateQuery + "?" + bm.EncodeURLParams()
	res, bs, err := a.doGet(ctx, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardTemplateQueryRsp{StatusCode: res.StatusCode}
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

// 会员卡模板修改 alipay.marketing.card.template.modify
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardTemplateModify(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardTemplateModifyRsp, err error) {
	err = bm.CheckEmptyError("request_id", "template_id", "template_style_info")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingCardTemplateModify, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardTemplateModifyRsp{StatusCode: res.StatusCode}
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

// 会员卡开卡表单模板配置 alipay.marketing.card.formtemplate.set
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardFormTemplateSet(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardFormTemplateSetRsp, err error) {
	err = bm.CheckEmptyError("template_id", "fields")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingCardFormTemplateSet, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardFormTemplateSetRsp{StatusCode: res.StatusCode}
	if res.StatusCode != http.StatusOK {
		if err = json.Unmarshal(bs, &aliRsp.ErrResponse); err != nil {
			return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
		}
		return aliRsp, nil
	}
	return aliRsp, a.autoVerifySignByCert(res, bs)
}

// 会员卡查询 alipay.marketing.card.query
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardQuery(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardQueryRsp, err error) {
	err = bm.CheckEmptyError("target_card_no_type", "target_card_no")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingCardQuery, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardQueryRsp{StatusCode: res.StatusCode}
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

// 会员卡更新 alipay.marketing.card.update
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardUpdate(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardUpdateRsp, err error) {
	err = bm.CheckEmptyError("target_card_no_type", "target_card_no", "occur_time", "card_info")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingCardUpdate, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardUpdateRsp{StatusCode: res.StatusCode}
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

// 会员卡删卡 alipay.marketing.card.delete
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardDelete(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardDeleteRsp, err error) {
	err = bm.CheckEmptyError("out_serial_no", "target_card_no", "target_card_no_type", "reason_code")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	uri := v3MarketingCardDelete + "?" + bm.EncodeURLParams()
	res, bs, err := a.doDelete(ctx, nil, uri, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardDeleteRsp{StatusCode: res.StatusCode}
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

// 会员卡消息通知 alipay.marketing.card.message.notify
// StatusCode = 200 is success
func (a *ClientV3) MarketingCardMessageNotify(ctx context.Context, bm gopay.BodyMap) (aliRsp *MarketingCardMessageNotifyRsp, err error) {
	err = bm.CheckEmptyError("target_card_no_type", "target_card_no", "occur_time")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	res, bs, err := a.doPost(ctx, bm, v3MarketingCardMessageNotify, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &MarketingCardMessageNotifyRsp{StatusCode: res.StatusCode}
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

// 上传门店照片和视频接口 alipay.offline.material.image.upload
// StatusCode = 200 is success
func (a *ClientV3) OfflineMaterialImageUpload(ctx context.Context, bm gopay.BodyMap) (aliRsp *OfflineMaterialImageUploadRsp, err error) {
	err = bm.CheckEmptyError("image_type", "image_name", "image_content")
	if err != nil {
		return nil, err
	}
	aat := bm.GetString(HeaderAppAuthToken)
	bm.Remove(HeaderAppAuthToken)
	res, bs, err := a.doProdPostFile(ctx, bm, v3OfflineMaterialImageUpload, aat)
	if err != nil {
		return nil, err
	}
	aliRsp = &OfflineMaterialImageUploadRsp{StatusCode: res.StatusCode}
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
