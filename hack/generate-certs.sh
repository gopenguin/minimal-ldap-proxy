#!/bin/bash

#
# a simple way to parse shell script arguments
#
# please edit and use to your hearts content
#


KEY_ALGORITHM="ecdsa"
KEY_PARAM="secp384r1"

function usage()
{
    echo "This script will generate a self signed certificate for testing"
    echo ""
    echo "$0"
    echo "    -h --help                         Print this help"
    echo "    --algorithm=$KEY_ALGORITHM        The algorithm used in the certificate: 'ecdsa' or 'rsa'"
    echo "    --param=$KEY_PARAM                Additional parameters to the algorithm, eg. ecdsa:secp384r1 or rsa:4096"
    echo ""
}

function setDefaults()
{
    case "$KEY_ALGORITHM" in
        ecdsa)
            if [ "$KEY_PARAM" == "" ]; then
                KEY_PARAM="secp384r1"
            fi
            ;;
        rsa)
            if [ "$KEY_PARAM" == "" ]; then
                KEY_PARAM="4096"
            fi
            ;;
    esac
}

while [ "$1" != "" ]; do
    PARAM=`echo $1 | awk -F= '{print $1}'`
    VALUE=`echo $1 | awk -F= '{print $2}'`
    case "$PARAM" in
        -h | --help)
            usage
            exit
            ;;
        --algorithm)
            KEY_ALGORITHM="$VALUE"
            ;;
        --param)
            KEY_PARAM="$VALUE"
            ;;
        *)
            echo "ERROR: unknown parameter \"$PARAM\""
            usage
            exit 1
            ;;
    esac
    shift
done

setDefaults

case "$KEY_ALGORITHM" in
    ecdsa)
        openssl req -x509 -nodes -newkey ec:<(openssl ecparam -name "$KEY_PARAM") -keyout server.key -out server.crt -days 3650
        ;;
    rsa)
        openssl req -x509 -nodes -newkey "rsa:$KEY_PARAM" -keyout server.key -out server.crt -days 3650
        ;;
    *)
        echo "ERROR: unsupported key algorithm."
        usage
        exit 1
        ;;
esac