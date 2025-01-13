package com.bukodi.playground.pqcbc;

import org.bouncycastle.jcajce.interfaces.MLDSAPrivateKey;
import org.bouncycastle.jcajce.interfaces.MLDSAPublicKey;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.jcajce.spec.MLKEMParameterSpec;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.pqc.jcajce.spec.DilithiumParameterSpec;
import org.bouncycastle.pqc.jcajce.spec.KyberParameterSpec;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.security.*;

public class TestBouncyCastlePQC {

    @BeforeClass
    public static void initBCPQC() {
        if (Security.getProvider("BC") == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }


    @Test
    public void testMLDSAApi() throws Exception {
        KeyPair kp = PQCUtil.generateKeyPair(MLDSAParameterSpec.ml_dsa_44);
        String privatePem = PQCUtil.exportKey((MLDSAPrivateKey) kp.getPrivate());
        System.out.println("Private key PEM: " + privatePem);

        String publicPem = PQCUtil.exportKey((MLDSAPublicKey) kp.getPublic());
        System.out.println("Public key PEM: " + publicPem);

        MLDSAPublicKey pubKey2 = PQCUtil.importPublicKey(publicPem);
        Assert.assertEquals("Public key", kp.getPublic(), pubKey2);

        MLDSAPrivateKey privKey2 = PQCUtil.importPrivateKey(privatePem);
        Assert.assertEquals("Private key", kp.getPrivate(), privKey2);
        MLDSAPublicKey pubKey3 = privKey2.getPublicKey();
        Assert.assertEquals("Public key", kp.getPublic(), pubKey3);

    }

    @Test
    public void testDilithium5Api() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(DilithiumParameterSpec.dilithium5.getName(), "BC");
        kpg.initialize(DilithiumParameterSpec.dilithium5, new SecureRandom());
        KeyPair kp = kpg.generateKeyPair();

        String privatePem = PQCUtil.exportKey( kp.getPrivate());
        System.out.println("Private key PEM: " + privatePem);

        String publicPem = PQCUtil.exportKey( kp.getPublic());
        System.out.println("Public key PEM: " + publicPem);

        PublicKey pubKey2 = PQCUtil.importPublicKey(publicPem);
        Assert.assertEquals("Public key", kp.getPublic(), pubKey2);

        PrivateKey privKey2 = PQCUtil.importPrivateKey(privatePem);
        Assert.assertEquals("Private key", kp.getPrivate(), privKey2);

    }

    @Test
    public void testMLKEMApi() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(MLKEMParameterSpec.ml_kem_512.getName(), "BC");
        kpg.initialize(MLKEMParameterSpec.ml_kem_512, new SecureRandom());
        KeyPair kp = kpg.generateKeyPair();
        String privatePem = PQCUtil.exportKey( kp.getPrivate());
        System.out.println("Private key PEM: " + privatePem);

        String publicPem = PQCUtil.exportKey( kp.getPublic());
        System.out.println("Public key PEM: " + publicPem);

        PublicKey pubKey2 = PQCUtil.importPublicKey(publicPem);
        Assert.assertEquals("Public key", kp.getPublic(), pubKey2);

        PrivateKey privKey2 = PQCUtil.importPrivateKey(privatePem);
        Assert.assertEquals("Private key", kp.getPrivate(), privKey2);
    }

    @Test
    public void testKyberApi() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(KyberParameterSpec.kyber512.getName(), "BC");
        kpg.initialize(KyberParameterSpec.kyber512, new SecureRandom());
        KeyPair kp = kpg.generateKeyPair();
        String privatePem = PQCUtil.exportKey( kp.getPrivate());
        System.out.println("Private key PEM: " + privatePem);

        String publicPem = PQCUtil.exportKey( kp.getPublic());
        System.out.println("Public key PEM: " + publicPem);

        PublicKey pubKey2 = PQCUtil.importPublicKey(publicPem);
        Assert.assertEquals("Public key", kp.getPublic(), pubKey2);

        PrivateKey privKey2 = PQCUtil.importPrivateKey(privatePem);
        Assert.assertEquals("Private key", kp.getPrivate(), privKey2);
    }

}
