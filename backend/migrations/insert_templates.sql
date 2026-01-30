USE dmh;

-- 插入经典模板
INSERT INTO poster_template_configs (id, name, preview_image, config, status) VALUES
(1, '经典模板', '/templates/classic.jpg', '{"width":750,"height":1334,"background":"#FFFFFF","elements":[{"type":"text","content":"{{campaignName}}","x":50,"y":100,"fontSize":32,"color":"#333333","fontWeight":"bold","maxWidth":650},{"type":"text","content":"{{campaignDescription}}","x":50,"y":160,"fontSize":18,"color":"#666666","maxWidth":650},{"type":"qrcode","content":"{{distributorLink}}","x":275,"y":1000,"size":200},{"type":"text","content":"扫码参与活动","x":375,"y":1220,"fontSize":16,"color":"#999999","align":"center"}]}', 'active');

-- 插入简约模板
INSERT INTO poster_template_configs (id, name, preview_image, config, status) VALUES
(2, '简约模板', '/templates/simple.jpg', '{"width":750,"height":1334,"background":"#F5F5F5","elements":[{"type":"text","content":"{{campaignName}}","x":375,"y":200,"fontSize":36,"color":"#000000","fontWeight":"bold","align":"center","maxWidth":700},{"type":"qrcode","content":"{{distributorLink}}","x":275,"y":500,"size":200},{"type":"text","content":"长按识别二维码","x":375,"y":750,"fontSize":18,"color":"#666666","align":"center"}]}', 'active');

-- 插入时尚模板
INSERT INTO poster_template_configs (id, name, preview_image, config, status) VALUES
(3, '时尚模板', '/templates/modern.jpg', '{"width":750,"height":1334,"background":"linear-gradient(135deg, #667eea 0%, #764ba2 100%)","elements":[{"type":"text","content":"{{campaignName}}","x":50,"y":150,"fontSize":40,"color":"#FFFFFF","fontWeight":"bold","maxWidth":650},{"type":"text","content":"{{campaignDescription}}","x":50,"y":230,"fontSize":20,"color":"#FFFFFF","maxWidth":650,"opacity":0.9},{"type":"rect","x":225,"y":900,"width":300,"height":300,"fill":"#FFFFFF","radius":20},{"type":"qrcode","content":"{{distributorLink}}","x":275,"y":950,"size":200},{"type":"text","content":"扫码立即参与","x":375,"y":1230,"fontSize":18,"color":"#FFFFFF","align":"center"}]}', 'active');
