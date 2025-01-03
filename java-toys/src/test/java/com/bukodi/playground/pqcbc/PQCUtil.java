package com.bukodi.playground.pqcbc;

import org.bouncycastle.jcajce.interfaces.MLDSAPrivateKey;
import org.bouncycastle.jcajce.interfaces.MLDSAPublicKey;
import org.bouncycastle.jcajce.provider.asymmetric.mldsa.MLDSAKeyPairGeneratorSpi;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.openssl.PEMWriter;
import org.bouncycastle.openssl.jcajce.JcaMiscPEMGenerator;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;
import org.bouncycastle.pqc.crypto.mldsa.MLDSAParameters;
import org.bouncycastle.pqc.crypto.mldsa.MLDSAPublicKeyParameters;
import org.bouncycastle.util.io.pem.PemGenerationException;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemObjectGenerator;
import org.bouncycastle.util.io.pem.PemReader;

import java.io.FileWriter;
import java.io.IOException;
import java.io.StringWriter;
import java.security.*;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;

public class PQCUtil {
    public static KeyPair generateKeyPair(MLDSAParameterSpec alg ) throws NoSuchAlgorithmException, NoSuchProviderException, InvalidAlgorithmParameterException {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(MLDSAParameterSpec.ml_dsa_44.getName(), "BC");
        kpg.initialize(MLDSAParameterSpec.ml_dsa_44, new SecureRandom());
        return kpg.generateKeyPair();
    }

    public static String exportPrivateKey( MLDSAPrivateKey mldsaPrivKey ) throws IOException {
        JcaMiscPEMGenerator pemGen = new JcaMiscPEMGenerator(mldsaPrivKey);
        StringWriter sw = new StringWriter();
        try( JcaPEMWriter pemWriter = new JcaPEMWriter(sw) ) {
            pemWriter.writeObject( pemGen);;
        };
        return sw.toString();
    }

    public static String exportPublicKey(MLDSAPublicKey mldsaPubKey ) throws IOException {
        JcaMiscPEMGenerator pemGen = new JcaMiscPEMGenerator(mldsaPubKey);
        StringWriter sw = new StringWriter();
        try( JcaPEMWriter pemWriter = new JcaPEMWriter(sw) ) {
            pemWriter.writeObject( pemGen);;
        };
        return sw.toString();
    }

    public static MLDSAPublicKey importPublicKey( String pubKeyPEM ) throws IOException, NoSuchAlgorithmException, NoSuchProviderException, InvalidKeySpecException {
        PemReader pemReader = new PemReader( new java.io.StringReader(pubKeyPEM) );
        PemObject pemObj = pemReader.readPemObject();
        byte[] asn1Bytes = pemObj.getContent();
        X509EncodedKeySpec keySpec = new X509EncodedKeySpec(asn1Bytes);
        KeyFactory keyFactory = KeyFactory.getInstance( MLDSAParameters.ml_dsa_44.getName(), "BC");
        PublicKey pubKey = keyFactory.generatePublic(keySpec);
        return (MLDSAPublicKey) pubKey;
    }

    public static MLDSAPrivateKey importPrivateKey( String privKeyPEM ) throws IOException, NoSuchAlgorithmException, NoSuchProviderException, InvalidKeySpecException {
        PemReader pemReader = new PemReader( new java.io.StringReader(privKeyPEM) );
        PemObject pemObj = pemReader.readPemObject();
        byte[] asn1Bytes = pemObj.getContent();
        PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(asn1Bytes);
        KeyFactory keyFactory = KeyFactory.getInstance( MLDSAParameters.ml_dsa_44.getName(), "BC");
        PrivateKey privKey = keyFactory.generatePrivate(keySpec);
        return (MLDSAPrivateKey)  privKey;
    }
}
