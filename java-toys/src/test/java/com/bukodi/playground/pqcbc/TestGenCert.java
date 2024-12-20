package com.bukodi.playground.pqcbc;

import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.DERSequence;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.*;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509CertificateHolder;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jcajce.spec.MLDSAParameterSpec;
import org.bouncycastle.jce.X509KeyUsage;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.jce.spec.ECNamedCurveGenParameterSpec;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.operator.jcajce.JcaContentVerifierProviderBuilder;
import org.junit.BeforeClass;
import org.junit.Test;

import java.math.BigInteger;
import java.nio.file.Files;
import java.security.*;
import java.security.cert.X509Certificate;
import java.util.Date;

public class TestGenCert {

    @BeforeClass
    public static void initBCPQC() {
        if (Security.getProvider("BC") == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    private static long ONE_YEAR = 365 * 24 * 60 * 60 * 1000L;

    @Test
    public void testGenDualCert() throws Exception {
        KeyPairGenerator kpGen = KeyPairGenerator.getInstance("ML-DSA", "BC");
        kpGen.initialize(MLDSAParameterSpec.ml_dsa_44, new SecureRandom());
        KeyPair kp = kpGen.generateKeyPair();
        PrivateKey privKey = kp.getPrivate();
        PublicKey pubKey = kp.getPublic();
        KeyPairGenerator ecKpGen = KeyPairGenerator.getInstance("EC", "BC");
        ecKpGen.initialize(new ECNamedCurveGenParameterSpec("P-256"), new SecureRandom());
        KeyPair ecKp = ecKpGen.generateKeyPair();
        PrivateKey ecPrivKey = ecKp.getPrivate();
        PublicKey ecPubKey = ecKp.getPublic();
        X500Name issuer = new X500Name("CN=ML-DSA ECDSA Alt Extension Certificate");
//
// create base certificate - version 3
//
        ContentSigner sigGen = new JcaContentSignerBuilder("SHA256withECDSA").setProvider("BC").build(ecPrivKey);
        ContentSigner altSigGen = new JcaContentSignerBuilder("ML-DSA-44").setProvider("BC").build(privKey);
        X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(
                issuer, BigInteger.valueOf(1),
                new Date(System.currentTimeMillis() - 50000),
                new Date(System.currentTimeMillis() + ONE_YEAR),
                issuer, ecPubKey)
                .addExtension(Extension.keyUsage, true,
                        new X509KeyUsage(X509KeyUsage.digitalSignature))
                .addExtension(Extension.extendedKeyUsage, true,
                        new DERSequence(KeyPurposeId.anyExtendedKeyUsage))
                .addExtension(new ASN1ObjectIdentifier("2.5.29.17"), true,
                        new GeneralNames(new GeneralName(GeneralName.rfc822Name, "test@test.test")))
                .addExtension(Extension.subjectAltPublicKeyInfo, false,
                        SubjectAltPublicKeyInfo.getInstance(kp.getPublic().getEncoded()));
        X509Certificate cert = new JcaX509CertificateConverter().setProvider("BC").getCertificate(certGen.build(sigGen, false,
                altSigGen));
// check validity and verify
        cert.checkValidity(new Date());
        cert.verify(cert.getPublicKey());
// create a certificate holder to allow checking of the altSignature.
        X509CertificateHolder certHolder = new JcaX509CertificateHolder(cert);
        SubjectPublicKeyInfo altPubKey =
                SubjectPublicKeyInfo.getInstance(certHolder.getExtension(Extension.subjectAltPublicKeyInfo).getParsedValue());
        if (certHolder.isAlternativeSignatureValid(new JcaContentVerifierProviderBuilder().setProvider("BC").build(altPubKey)))

        {
            System.out.println("alternative signature verified on certificate");
        } else {
            System.out.println("alternative signature verification failed on certificate");
        }

        byte[] derCer = certHolder.getEncoded();
        Files.write(new java.io.File("/tmp/dualcert.cer").toPath(), derCer);

    }
    @Test
    public void testGenCert() throws Exception {
// generate an ML-DSA-44 key pair
        KeyPairGenerator kpGen = KeyPairGenerator.getInstance("ML-DSA", "BC");
        kpGen.initialize(MLDSAParameterSpec.ml_dsa_44, new SecureRandom());
        KeyPair kp = kpGen.generateKeyPair();
        PrivateKey privKey = kp.getPrivate();
        PublicKey pubKey = kp.getPublic();
        X500Name issuer = new X500Name("CN=ML-DSA Certificate");

//
// create base certificate - version 3
//
        ContentSigner sigGen = new JcaContentSignerBuilder("ML-DSA-44").setProvider("BC").build(privKey);
        X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(
                issuer, BigInteger.valueOf(1),
                new Date(System.currentTimeMillis() - 50000),
                new Date(System.currentTimeMillis() + ONE_YEAR),
                issuer, pubKey)
                .addExtension(Extension.keyUsage, true,
                        new X509KeyUsage(X509KeyUsage.encipherOnly))
                .addExtension(Extension.extendedKeyUsage, true,
                        new DERSequence(KeyPurposeId.anyExtendedKeyUsage))
                .addExtension(Extension.subjectAlternativeName, true,
                        new GeneralNames(new GeneralName(GeneralName.rfc822Name, "test@test.test")));
        X509Certificate cert = new JcaX509CertificateConverter().setProvider("BC").getCertificate(certGen.build(sigGen));
//
// check validity
//
        cert.checkValidity(new Date());
        cert.verify(cert.getPublicKey());
        System.out.println("ML-DSA certificate verified");
    }
}
