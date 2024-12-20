package com.bukodi.playground.pqcbc;

import org.bouncycastle.jcajce.provider.asymmetric.mldsa.MLDSAKeyPairGeneratorSpi;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.BeforeClass;
import org.junit.Test;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.Security;
import java.util.Arrays;

public class TestBouncyCastlePQC {

    @BeforeClass
    public static void initBCPQC() {
        if (Security.getProvider("BC") == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    @Test
    public void testMLDSA() throws Exception {
        //MLDSAKeyPairGenerator kpg = new MLDSAKeyPairGenerator();
        //KeyGenerationParameters kgParams = new MLDSAKeyGenerationParameters(256, 512);
        //kpg.init( kgParams );
        Arrays.asList(Security.getProviders()).forEach(p -> System.out.println(p.getName()));
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(MLDSAParameterSpec.ml_dsa_44.getName(), "BC");
        if (kpg instanceof MLDSAKeyPairGeneratorSpi) {
            System.out.println("MLDSAKeyPairGenerator");
        }
        //kpg.initialize(MLDSASpec.mldsa512, new SecureRandom());
        KeyPair keyPair = kpg.generateKeyPair();
        //
        // // get private and public key
        // PrivateKey privateKey = keyPair.getPrivate();
        // BCPQCMultiLayerPrivateKey mldsaPrivateKey = (BCPQCMultiLayerPrivateKey)privateKey;
        // //PrivateKeyInfo privKeyInfo = PrivateKeyInfoFactory.createPrivateKeyInfo(mldsaPrivateKey.getParameterSpec());
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

}
