package com.bukodi.playground.pqcbc;

import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.SubjectPublicKeyInfo;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.openssl.jcajce.JcaMiscPEMGenerator;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;
import org.bouncycastle.util.io.pem.PemReader;

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

    public static String exportKey( Key key ) throws IOException {
        StringWriter sw = new StringWriter();
        try( JcaPEMWriter pemWriter = new JcaPEMWriter(sw) ) {
            pemWriter.writeObject( new JcaMiscPEMGenerator(key));;
        };
        return sw.toString();
    }

    @SuppressWarnings("unchecked")
    public static <T extends PublicKey> T importPublicKey(String pubKeyPEM ) throws IOException, NoSuchAlgorithmException, NoSuchProviderException, InvalidKeySpecException {
        byte[] asn1Bytes = (new PemReader( new java.io.StringReader(pubKeyPEM) )).readPemObject().getContent();

        X509EncodedKeySpec keySpec = new X509EncodedKeySpec(asn1Bytes);
        AlgorithmIdentifier alg = SubjectPublicKeyInfo.getInstance(keySpec.getEncoded()).getAlgorithm();

        PublicKey pubKey = KeyFactory.getInstance( alg.getAlgorithm().getId()).generatePublic(keySpec);
        return (T) pubKey;
    }


    @SuppressWarnings("unchecked")
    public static <T extends PrivateKey> T importPrivateKey(String privKeyPEM ) throws IOException, NoSuchAlgorithmException, NoSuchProviderException, InvalidKeySpecException {
        byte[] asn1Bytes = (new PemReader( new java.io.StringReader(privKeyPEM) )).readPemObject().getContent();

        PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(asn1Bytes);
        ASN1ObjectIdentifier algOid = PrivateKeyInfo.getInstance(keySpec.getEncoded()).getPrivateKeyAlgorithm().getAlgorithm();
        PrivateKey privKey = KeyFactory.getInstance(algOid.getId()).generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        return (T)  privKey;
    }

}
