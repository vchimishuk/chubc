PORTNAME=	chubc
PORTVERSION=	0.0.1
CATEGORIES=	audio

MAINTAINER=	vchimishuk@yandex.ru
COMMENT=	non-interactive console client for Chub audio player
WWW=		https://github.com/vchimishuk/chubc

LICENSE=	GPLv3+

FETCH_DEPENDS=	git:devel/git
USES=		go:1.21,modules

USE_GITHUB=	yes
GH_ACCOUNT=	vchimishuk
GH_PROJECT=	chubc
GH_TAGNAME=	7fcf3b3
# $ make gomod-vendor
GH_TUPLE=	vchimishuk:chubby:edf89263defb:vchimishuk_chubby/vendor/github.com/vchimishuk/chubby \
		vchimishuk:opt:e09d5e511662:vchimishuk_opt/vendor/github.com/vchimishuk/opt

PLIST_FILES=	bin/chubc \
		share/man/man1/chubc.1.gz

post-install:
	${INSTALL_MAN} ${WRKSRC}/chubc.1 ${STAGEDIR}${PREFIX}/share/man/man1

.include <bsd.port.mk>
