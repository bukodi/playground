package com.bukodi.playground.pqcbc;

import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.crypto.params.AsymmetricKeyParameter;
import org.bouncycastle.crypto.util.PrivateKeyInfoFactory;
import org.bouncycastle.jcajce.interfaces.MLDSAPrivateKey;
import org.bouncycastle.jcajce.interfaces.MLDSAPublicKey;
import org.bouncycastle.jcajce.provider.asymmetric.mldsa.BCMLDSAPrivateKey;
import org.bouncycastle.jcajce.provider.asymmetric.mldsa.MLDSAKeyPairGeneratorSpi;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.pqc.crypto.mldsa.MLDSAParameters;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.security.*;
import java.util.Arrays;

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
        String privatePem = PQCUtil.exportPrivateKey((MLDSAPrivateKey) kp.getPrivate());
        System.out.println("Private key PEM: " + privatePem);

        String publicPem = PQCUtil.exportPublicKey((MLDSAPublicKey) kp.getPublic());
        System.out.println("Public key PEM: " + publicPem);

        MLDSAPublicKey pubKey2 = PQCUtil.importPublicKey(publicPem);
        Assert.assertEquals("Public key", kp.getPublic(), pubKey2);

        MLDSAPrivateKey privKey2 = PQCUtil.importPrivateKey(privatePem);
        Assert.assertEquals("Private key", kp.getPrivate(), privKey2);
        MLDSAPublicKey pubKey3 = privKey2.getPublicKey();
        Assert.assertEquals("Public key", kp.getPublic(), pubKey3);

    }


        @Test
    public void testMLDSA() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(MLDSAParameterSpec.ml_dsa_44.getName(), "BC");
        if (kpg instanceof MLDSAKeyPairGeneratorSpi) {
            System.out.println("MLDSAKeyPairGenerator");
        }
        kpg.initialize(MLDSAParameterSpec.ml_dsa_44, new SecureRandom());
        KeyPair keyPair = kpg.generateKeyPair();

        // get private and public key
        PrivateKey privateKey = keyPair.getPrivate();
        System.out.println(privateKey.getAlgorithm());
        MLDSAPrivateKey mldsaPrivateKey = (MLDSAPrivateKey)privateKey;
        MLDSAParameterSpec pramSpec = mldsaPrivateKey.getParameterSpec();

        PublicKey publicKey = keyPair.getPublic();
        System.out.println(publicKey.getAlgorithm());
        MLDSAPublicKey mldsaPublicKey = (MLDSAPublicKey) publicKey;
        MLDSAParameterSpec pramSpecPub = mldsaPublicKey.getParameterSpec();

        KeyFactory keyFactory = KeyFactory.getInstance(MLDSAParameters.ml_dsa_44.getName(), "BC");
        System.out.println(keyFactory.getClass().getName());

        // //MLDSAPrivateKeyParameters mldsaPrivateKeyParams = (MLDSAPrivateKeyParameters)privKeyInfo;
        // PublicKey publicKey = keyPair.getPublic();
        // BCPQCMultiLayerPublicKey mldsaPublicKey = (BCPQCMultiLayerPublicKey) publicKey;
        //
        // // storing the key as byte array
        // byte[] privateKeyByte = mldsaPrivateKey.getEncoded();
        // byte[] publicKeyByte = publicKey.getEncoded();
        // System.out.printf( "Private key %d bytes: %s\n", privateKeyByte.length, byteArrayToHex(privateKeyByte));
        // System.out.printf( "Public key %d bytes: %s\n", publicKeyByte.length, byteArrayToHex(publicKeyByte));
    }

    public KeyPair generateKeyPair( MLDSAParameterSpec alg ) {

        return null;
    }
}
