#!/bin/bash -x
die() {
    if [[ $1 -eq 0 ]]; then
        rm -rf "${TMPDIR}"
    else
        [ "${TMPDIR}" = "" ] || echo "Temporary data stored at '${TMPDIR}'"
    fi
    echo "$2"
    exit $1
}

pwd
TIMESTAMP=$( git log -1 -s --format=%ct HEAD )
VERSION=$( date +'%Y.%m.%d' -d @${TIMESTAMP} )
RELEASE=${TIMESTAMP}.$( git rev-parse --short HEAD )

egrep -q '^[0-9]+\.[a-z0-9]+$' <<< ${RELEASE} || {
	echo "Revision: ${RELEASE}";
	echo "Version: ${VERSION}";
	echo
	echo "Known tags:"
	git tag
	echo;
	die 1 "Can't get latest version from git";
}

TMPDIR=$(mktemp -d)
echo ${VERSION}-${RELEASE}
make || die 1 "Can't build package"
mkdir -p "${TMPDIR}/usr/bin" "${TMPDIR}/usr/share" || die 1 "Can't create install dir"
cp ./carbonapi "${TMPDIR}/usr/bin" || die 1 "Can't install package"
cp ./cmd/carbonapi/*.example.* "${TMPDIR}/usr/share" || die 1 "Can't install share"
mkdir -p "${TMPDIR}"/etc/systemd/system/
mkdir -p "${TMPDIR}"/etc/sysconfig/
cp ./contrib/carbonapi/rhel/carbonapi.service "${TMPDIR}"/etc/systemd/system/
cp ./contrib/carbonapi/common/carbonapi.env "${TMPDIR}"/etc/sysconfig/carbonapi

fpm -s dir -t rpm -n carbonapi -v ${VERSION} -C ${TMPDIR} \
    --iteration ${RELEASE} \
    -p carbonapi-VERSION-ITERATION.ARCH.rpm \
    -d "cairo" \
    --after-install contrib/carbonapi/fpm/systemd-reload.sh \
    --description "carbonapi: replacement graphite API server" \
    --license BSD-2 \
    --url "https://github.com/bookingcom/carbonapi" \
    "${@}" \
    etc usr/bin usr/share || die 1 "Can't create package!"

die 0 "Success"
