package com.bukodi.playground.pqcbc;

import org.bouncycastle.jcajce.interfaces.MLDSAPrivateKey;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.openssl.jcajce.JcaMiscPEMGenerator;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;
import org.bouncycastle.pqc.crypto.mldsa.MLDSAParameters;
import org.bouncycastle.util.io.pem.PemReader;
import org.bouncycastle.util.test.FixedSecureRandom;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.StringWriter;
import java.security.*;
import java.security.spec.PKCS8EncodedKeySpec;

public class TestPQCPrivateKeyASN1Bug {

    @BeforeClass
    public static void initBCPQC() {
        if (Security.getProvider("BC") == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    @Test
    public void exportMLDSAPrivKey() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance(MLDSAParameterSpec.ml_dsa_44.getName(), "BC");
        kpg.initialize(MLDSAParameterSpec.ml_dsa_44, new FixedSecureRandom( new byte[] { 
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
                0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
        }));
        KeyPair kp = kpg.generateKeyPair();

        MLDSAPrivateKey mldsaPrivKey = (MLDSAPrivateKey) kp.getPrivate();
        StringWriter sw = new StringWriter();
        try( JcaPEMWriter pemWriter = new JcaPEMWriter(sw) ) {
            pemWriter.writeObject( new JcaMiscPEMGenerator(mldsaPrivKey));;
        };
        System.out.println("Private key PEM: \n" + sw);
    }

    private static String bcGeneratedMLDSA44 = "-----BEGIN PRIVATE KEY-----\n" +
            "MDQCAQAwCwYJYIZIAWUDBAMRBCIEIAABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZ\n" +
            "GhscHR4f\n" +
            "-----END PRIVATE KEY-----";


    // See https://datatracker.ietf.org/doc/draft-ietf-lamps-dilithium-certificates/ Appecdix C.1.
    private static String ietfExampleMLDSA44 = "-----BEGIN PRIVATE KEY-----\n" +
            "MDICAQAwCwYJYIZIAWUDBAMRBCAAAQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRob\n" +
            "HB0eHw==\n" +
            "-----END PRIVATE KEY-----";

    @Test
    public void importMLDSAPrivKeyBc() throws Exception {
        String pem = bcGeneratedMLDSA44;
        byte[] asn1Bytes = (new PemReader( new java.io.StringReader(pem) )).readPemObject().getContent();
        PrivateKey privKey = KeyFactory.getInstance( MLDSAParameters.ml_dsa_44.getName(), "BC").generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        System.out.println("Private key: " + privKey);
    }

    @Test
    public void importMLDSAPrivKeyIetf() throws Exception {
        String pem = ietfExampleMLDSA44;
        byte[] asn1Bytes = (new PemReader( new java.io.StringReader(pem) )).readPemObject().getContent();
        PrivateKey privKey = KeyFactory.getInstance( MLDSAParameters.ml_dsa_44.getName(), "BC").generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        System.out.println("Private key: " + privKey);
    }


}

