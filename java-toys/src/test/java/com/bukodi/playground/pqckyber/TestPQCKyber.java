package com.bukodi.playground.pqckyber;

import org.bouncycastle.pqc.jcajce.interfaces.KyberPrivateKey;
import org.bouncycastle.pqc.jcajce.provider.BouncyCastlePQCProvider;
import org.bouncycastle.pqc.jcajce.spec.KyberParameterSpec;
import org.junit.BeforeClass;
import org.junit.Test;

import java.security.*;

public class TestPQCKyber {

    @BeforeClass
    public static void initBCPQC() {
        if (Security.getProvider("BCPQC") == null) {
            Security.addProvider(new BouncyCastlePQCProvider());
        }
    }

    public static String byteArrayToHex(byte[] a) {
        StringBuilder sb = new StringBuilder(a.length * 2);
        for(byte b: a)
            sb.append(String.format("%02x", b));
        return sb.toString();
    }

    @Test
    public void testPQCKyber() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("KYBER", "BCPQC");
        kpg.initialize(KyberParameterSpec.kyber768, new SecureRandom());
        KeyPair keyPair = kpg.generateKeyPair();

        // get private and public key
        PrivateKey privateKey = keyPair.getPrivate();
        KyberPrivateKey kyberPrivateKey = (KyberPrivateKey)privateKey;
        //PrivateKeyInfo privKeyInfo = PrivateKeyInfoFactory.createPrivateKeyInfo(privateKey);
        //KyberPrivateKeyParameters kyberPrivateKeyParams = (KyberPrivateKeyParameters)privKeyInfo;
        PublicKey publicKey = keyPair.getPublic();

        // storing the key as byte array
        byte[] privateKeyByte = privateKey.getEncoded();
        byte[] publicKeyByte = publicKey.getEncoded();
        System.out.printf( "Private key %d bytes: %s\n", privateKeyByte.length, byteArrayToHex(privateKeyByte));
        System.out.printf( "Public key %d bytes: %s\n", publicKeyByte.length, byteArrayToHex(publicKeyByte));

    }
}
