import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.SubjectPublicKeyInfo;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.openssl.jcajce.JcaMiscPEMGenerator;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;
import org.bouncycastle.util.io.pem.PemReader;

import java.io.IOException;
import java.io.StringWriter;
import java.security.*;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;

public class BugReproducePQCPrivateKeyASN1 {

    private static String bcGeneratedMLDSA44 = "-----BEGIN PRIVATE KEY-----\n" +
            "MDQCAQAwCwYJYIZIAWUDBAMRBCIEIAABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZ\n" +
            "GhscHR4f\n" +
            "-----END PRIVATE KEY-----";


    // See https://datatracker.ietf.org/doc/draft-ietf-lamps-dilithium-certificates/06/ Appecdix C.1.
    private static String ietfExampleMLDSA44 = "-----BEGIN PRIVATE KEY-----\n" +
            "MDICAQAwCwYJYIZIAWUDBAMRBCAAAQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRob\n" +
            "HB0eHw==\n" +
            "-----END PRIVATE KEY-----";


    // See https://datatracker.ietf.org/doc/draft-ietf-lamps-kyber-certificates/07/ Appecdix C.1.
    private static String ietfExampleMLKEM512 = "-----BEGIN PRIVATE KEY-----\n" +
            "MDICAQAwCwYJYIZIAWUDBAMRBCAAAQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRob\n" +
            "HB0eHw==\n" +
            "-----END PRIVATE KEY-----";

    public static void main(String[] args) throws Exception {
        Security.addProvider(new BouncyCastleProvider());

        System.out.println("BC generated key:\n" + bcGeneratedMLDSA44);
        System.out.println();
        System.out.println("IETF example key:\n" + ietfExampleMLDSA44);




        String pem = ietfExampleMLDSA44;
        byte[] asn1Bytes = (new PemReader( new java.io.StringReader(pem) )).readPemObject().getContent();
        PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(asn1Bytes);
        PrivateKeyInfo privKeyInfo = PrivateKeyInfo.getInstance(keySpec.getEncoded());
        ASN1ObjectIdentifier algOid = privKeyInfo.getPrivateKeyAlgorithm().getAlgorithm();
        PrivateKey privKey = KeyFactory.getInstance(algOid.getId(), "BC").generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        //PrivateKey privKey = KeyFactory.getInstance( MLDSAParameters.ml_dsa_44.getName(), "BC").generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        System.out.println("Private key: " + privKey);
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
