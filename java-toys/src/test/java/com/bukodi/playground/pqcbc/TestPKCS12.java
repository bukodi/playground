package com.bukodi.playground.pqcbc;

import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.pqc.crypto.mldsa.MLDSAParameters;
import org.bouncycastle.util.io.pem.PemReader;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.FileOutputStream;
import java.security.KeyFactory;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.Security;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;
import java.security.spec.PKCS8EncodedKeySpec;

public class TestPKCS12 {

    @BeforeClass
    public static void initBCPQC() {
        if (Security.getProvider("BC") == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }


    @Test
    public void generatePKCS12() throws Exception {
        // Load private key
        byte[] asn1Bytes = (new PemReader(new java.io.StringReader( bcGeneratedMLDSA44))).readPemObject().getContent();
        System.out.printf("asn1Bytes: %d\n", asn1Bytes.length);
        KeyFactory kpg = KeyFactory.getInstance(MLDSAParameters.ml_dsa_44.getName(), "BC");
        PrivateKey privKey = kpg.generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        System.out.println("Private key: " + privKey);

        // Load certificate
        byte[] certBytes = (new PemReader(new java.io.StringReader(ML_DSA_44_crt))).readPemObject().getContent();
        CertificateFactory certFactory = CertificateFactory.getInstance("X.509", "BC");
        Certificate certificate = certFactory.generateCertificate(new java.io.ByteArrayInputStream(certBytes));
        System.out.println("Certificate: " + certificate);

        // Create a KeyStore of type PKCS12
        KeyStore keyStore = KeyStore.getInstance("PKCS12", "BC");
        keyStore.load(null, null);

        // Store the private key and certificate chain in the KeyStore
        keyStore.setKeyEntry("alias", privKey, "Passw0rd".toCharArray(), new Certificate[]{certificate});

        // Write the KeyStore to a file
        try (FileOutputStream fos = new FileOutputStream("/tmp/keystore.p12")) {
            keyStore.store(fos, "Passw0rd".toCharArray());
        }

        System.out.println("PKCS12 file generated successfully.");

    }

    @Test
    public void generatePKCS12_EC() throws Exception {
        String privKeyTxt = "-----BEGIN PRIVATE KEY-----\n" +
                "MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgs1evReNiR1/S+91j\n" +
                "KHTwdxtJFtpMVQSuZG2EqkvjZY+gCgYIKoZIzj0DAQehRANCAASVmXNKeIzAK3is\n" +
                "Ic3CHJgK7zp8tgWLazNJrKAAmCVR0FxQrEPqOd+03BkwVRcQdMermoue3Ay/SoKZ\n" +
                "Fke0xdbs\n" +
                "-----END PRIVATE KEY-----    ";

        String cicaCert = "-----BEGIN CERTIFICATE-----\n" +
                "MIIDkTCCAXmgAwIBAgIUclS1B7UGvgUNjNMMBAQdID2jg+4wDQYJKoZIhvcNAQEL\n" +
                "BQAwOzEVMBMGA1UEAwwMTWFuYWdlbWVudENBMRUwEwYDVQQKDAxFSkJDQSBTYW1w\n" +
                "bGUxCzAJBgNVBAYTAlNFMB4XDTIzMTAwNjExMjcyNFoXDTI1MTAwNTExMjcyM1ow\n" +
                "FDESMBAGA1UEAwwJQ2ljYSBNaWNhMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\n" +
                "lZlzSniMwCt4rCHNwhyYCu86fLYFi2szSaygAJglUdBcUKxD6jnftNwZMFUXEHTH\n" +
                "q5qLntwMv0qCmRZHtMXW7KN/MH0wDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBRN\n" +
                "F+Zycm+de8qzbboh67phsFjGaTAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUH\n" +
                "AwQwHQYDVR0OBBYEFNidaA79z72ds1WDD2J9GI/m9QA/MA4GA1UdDwEB/wQEAwIF\n" +
                "4DANBgkqhkiG9w0BAQsFAAOCAgEA0ZzWHvZ9RjjM2EE54dZVAqcq5qD8Kg1Bm1uU\n" +
                "6TeBmTC/EzwaRy2F/CjFT6ego03tKHgAV8tAAjE1UteWqOj+gLO1Sld4dwN3mlS3\n" +
                "B76aRHmuOzD0etZx8Yi8td6Ja5PJcdn44GSU7jZMJ7+SoOqa373DQj6HHX+YjFnz\n" +
                "SbIvFCNrWJJxkKNDuVbMaafJStM4gunnD0ZLtPiD+NVTMzU2idafvGOlKAne1gkW\n" +
                "LfAbleqXblXOS0U7oCEIuYwdUXBa/W54Fx/MOuoVRrOQQ2Q2j1xRITm6vsrkNKP5\n" +
                "7EQISXo2XuiD5zBwmeBPKzs9qN7VNDhCqfddsq21pJSnXwwj1rKOEk0vZHzGcGeD\n" +
                "B0hjHQ8QrF8Jvh1aT0+tOgPeaPY4Bzw/ojTJBM97jqzHaY2UwZ+VlHdRTwAz5HIx\n" +
                "Yaot/vRNY7090ktrtLuzKjMcTFBXkxykes/YWwC9zxQRbz9jsM7cbBIaOiNtOQEt\n" +
                "vXkNO4lVhes5ckPdhWkpQUhnBLR5vW7cmIYHMvhQZ26kwUCBnqXHlmmiwEbQznb7\n" +
                "ANYmCbig7p/Tkbl1fQiiQNiDb6K+OoWBoqJS/easolKZfrh1Q6iAxwGQdHr3Lfj2\n" +
                "qh5JTEOUvmnIA/OW1z/GkkZCCmV5A/S5LxL9t+57VNSG1YIyq1IHFVTD7sixJqUD\n" +
                "q9Ndni4=\n" +
                "-----END CERTIFICATE-----";

        // Load private key
        byte[] asn1Bytes = (new PemReader(new java.io.StringReader( privKeyTxt))).readPemObject().getContent();
        KeyFactory kpg = KeyFactory.getInstance( "ECDSA", "BC");
        PrivateKey privKey = kpg.generatePrivate(new PKCS8EncodedKeySpec(asn1Bytes));
        System.out.println("Private key: " + privKey);

        // Load certificate
        byte[] certBytes = (new PemReader(new java.io.StringReader(cicaCert))).readPemObject().getContent();
        CertificateFactory certFactory = CertificateFactory.getInstance("X.509", "BC");
        Certificate certificate = certFactory.generateCertificate(new java.io.ByteArrayInputStream(certBytes));
        System.out.println("Certificate: " + certificate);

        // Create a KeyStore of type PKCS12
        KeyStore keyStore = KeyStore.getInstance("PKCS12", "BC");
        keyStore.load(null, null);

        // Store the private key and certificate chain in the KeyStore
        keyStore.setKeyEntry("cicaAlias", privKey, "Passw0rd".toCharArray(), new Certificate[]{certificate});

        // Write the KeyStore to a file
        try (FileOutputStream fos = new FileOutputStream("/tmp/keystore_cica.p12")) {
            keyStore.store(fos, "Passw0rd".toCharArray());
        }

        System.out.println("PKCS12 Cica file generated successfully.");

    }

    final static String ML_DSA_44_crt = "-----BEGIN CERTIFICATE-----\n" +
            "MIIPlDCCBgqgAwIBAgIUFZ/+byL9XMQsUk32/V4o0N44804wCwYJYIZIAWUDBAMR\n" +
            "MCIxDTALBgNVBAoTBElFVEYxETAPBgNVBAMTCExBTVBTIFdHMB4XDTIwMDIwMzA0\n" +
            "MzIxMFoXDTQwMDEyOTA0MzIxMFowIjENMAsGA1UEChMESUVURjERMA8GA1UEAxMI\n" +
            "TEFNUFMgV0cwggUyMAsGCWCGSAFlAwQDEQOCBSEA17K0clSq4NtF55MNSpjSyX2P\n" +
            "E5fReJ2voXAksxbpvslPyZRtQvGbeadBO7qjPnFJy0LtURVpOsBB+suYit61/g4d\n" +
            "hjEYSZW1ksOX0ilOLhT5CqQUujgmiZrEP0zMrLwm6agyuVEY1ctDPL75ZgsAE44I\n" +
            "F/YediyidMNq1VTrIqrBFi5KsBrLoeOMTv2PgLZbMz0PcuVd/nHOnB67mInnxWEG\n" +
            "wP1zgDoq7P6v3teqPLLO2lTRK9jNNqeM+XWUO0er0l6ICsRS5XQu0ejRqCr6huWQ\n" +
            "x1jBWuTShA2SvKGlCQ9ASWWX/KfYuVE/GhvabpUKqpjeRnUH1KT1pPBZkhZYLDVy\n" +
            "9i7aiQWrNYFnDEoCd3oz4Mpylf2PT/bRoKOnaD1l9fX3/GDaAj6CbF+SFEwC99G6\n" +
            "EHWYdVPqk2f8122ZC3+pnNRa/biDbUPkWfUYffBYR5cJoB6mg1k1+nBGCZDNPcG6\n" +
            "QBupS6sd3kGsZ6szGdysoGBI1MTu8n7hOpwX0FOPQw8tZC3CQVZg3niHfY2KvHJS\n" +
            "OXjAQuQoX0MZhGxEEmJCl2hEwQ5Va6IVtacZ5Z0MayqW05hZBx/cws3nUkp77a5U\n" +
            "6FsxjoVOj+Ky8+36yXGRKCcKr9HlBEw6T9r9n/MfkHhLjo5FlhRKDa9YZRHT2ZYr\n" +
            "nqla8Ze05fxg8rHtFd46W+9fib3HnZEFHZsoFudPpUUx79wcvnTUSIV/R2vNWPIc\n" +
            "C2U7O3ak4HamVZowJxhVXMY/dIWaq6uSXwI4YcqM0Pe62yhx9n1VMm10URNa1F9K\n" +
            "G6aRGPuyyKMO7JOS7z+XcGbJrdXHEMxkexUU0hfZWMcBfD6Q/SDATmdLkEhuk3Cj\n" +
            "GgAdMvRzl55JBnSefkd/oLdFCPil8jeDErg8Jb04jKCw//dHi69CtxZn7arJfEax\n" +
            "KWQ+WG5bBVoMIRlG1PNuZ1vtWGD6BCoxXZgmFk1qkjfDWl+/SVSQpb1N8ki5XEqu\n" +
            "d4S2BWcxZqxCRbW0sIKgnpMj5i8geMW3Z4NEbe/XNq06NwLUmwiYRJAKYYMzl7xE\n" +
            "GbMNepegs4fBkRR0xNQbU+Mql3rLbw6nXbZbs55Z5wHnaVfe9vLURVnDGncSK1IE\n" +
            "47XCGfFoixTtC8C4AbPm6C3NQ+nA6fQXRM2YFb0byIINi7Ej8E+s0bG2hd1aKxuN\n" +
            "u/PtkzZw8JWhgLTxktCLELj6u9/MKyRRjjLuoKXgyQTKhEeACD87DNLQuLavZ7w1\n" +
            "W5SUAl3HsKePqA46Lb/rUTKIUdYHgZjpSTZRrnh+wCUfkiujDp9R32Km1yeEzz3S\n" +
            "BTkxdt+jJKUSvZSXCjbdNKUUqGeR8Os28BRbCatkZRtKAxOymWEaKhxIiRYnWYdo\n" +
            "oxFAYLpEQ0ht9RUioc6IswmFwhb45u0XjdVnswSg1Mr7qIKig0LxepqiauWNtjAI\n" +
            "PSw1j99WbD9dYqQoVnvJ6ozpXKoPNUdLC/qPM5olCrTfzyCDvo7vvBBV4Y/hU3Du\n" +
            "yyYFZtg/8GshGq7EPKKbVMzQD4gVokZe8LRlFcx+QfMSTwnv/3OTCatYspoUWaAL\n" +
            "zlA46TjJZ49y6w5O5f2q5m2fhXP8l/xCtJWfS/i2HXhDPoawM11ukZHE2L9IezkF\n" +
            "wQjP1qwksM633LfPUfhNDtaHuV6uscUzwG8NlwI9kqcIJYN7Wbpst9TlawqHwgOG\n" +
            "KujzFbpZJejt76Z5NpoiAnZhUfFqll+fgeznbMBwtVhp5NuXhM8FyDCzJCyDEqNC\n" +
            "MEAwDgYDVR0PAQH/BAQDAgGGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFDKa\n" +
            "B7H6u0j1KjCfEaGJj4SOIyL/MAsGCWCGSAFlAwQDEQOCCXUAZ6iVH8MI4S9oZ2Ef\n" +
            "3CVL9Ly1FPf18v3rcvqOGgMAYWd7hM0nVZfYMVQZWWaxQWcMsOiBE0YNl4oaejiV\n" +
            "wRykGZV3XAnWTd60e8h8TovxyTJ/xK/Vw3hlU+F9YpsPJxQnZUgUMrXnzNC6YeUc\n" +
            "rT3Y+Vk4wjXr7O6vixauM2bzAMU1jse+nrI6HqGj2lhoZwTwSD+Wim5LH4lnCgE0\n" +
            "s2oY1scn3JsCexJ5R5OkjHq2bt9XrBgRORTADQoRtlplL0d3Eze/dDZm/Klby9OR\n" +
            "Ia4HUL7FWtWoy86Y5TiuUjlH1pKZdjMPyj/JXAHRQDtJ5cuoGBL0NlDdATEJNCee\n" +
            "zQfMqzTCyjCn091QkuFjDhQjzJ+sQ6G02w49lw8Kpm1ASuh7BLTPcuz7Z+rLpNjN\n" +
            "jmW67rR6+hHMK474mSKIZnuO3vVKnidntjLhSYc1soxvYPCLWWnl4m3XyjlrnlzD\n" +
            "4Soec2I2AjKNZKCO9KKa81cRzIcNJjc7sbnrLv/hKXNUTESn4s3yAyRPU7N6bVIy\n" +
            "N9ifBvb1U07WMRPI8A7/f9zVCaLYx87ym9P7GGpMjDYrPUQpOaKQdu4ycWuPrlEA\n" +
            "2BoHIVzbHHm9373BT1LjcxjR5SbbhNFg+42hwG284VlVzcLW/XiipaWN8jnONmxt\n" +
            "kLMui9R/wf0TCehilMDDtRznfm37b2ci5o9MP/LrTDRpMVBudDuwIZmLgPQ/bj08\n" +
            "n+VHd8D2WADpR/kEMpDhSwG2P44mwwE4CUKGbHS0qQLOSRwMlQVEzwxpOOrLMusw\n" +
            "JmzoLE0KNsUR6o/3xAlUmjqCZMqYPYxtXgNfJEJDp3V1iqyZK1iES3EQ0/h8m7oZ\n" +
            "3YqNKrEpTgVV7EmVpUjcVszjWgXcSKynVVsWQd3j0Zf83zXRLwmq8+anJ3XNGCSa\n" +
            "IecO2sZxDbaiHhwFYRkt0BGRM2QM//IPMYeXhRa/1svmbOEHGxJG9LqTffkBs+01\n" +
            "Bp7r3/9lRZ+5t3eukpinpJrCT0AgeV3l3ujbzyCiQbboFDaPS4+kKvi+iS2eHjiu\n" +
            "S/WkfP1Go5jksxhkceJFNPsTmGCyXGPy2/haU9hkiMg9/wmuIKm/gxRfIBh/DoIr\n" +
            "1HWZjTuWcBGWTu2NuXeAVO/MbMtpB0u6mWYktHQcVxA2LenU+N5LEPbbHp+AmPQC\n" +
            "RZPqBziTyx/nuVnFD+/EAbPKzeqMKhcTW6nfkKt/Md4zmi1vhWxx7c+wDlo9cyAf\n" +
            "vsS0p5uXKK1wzaC4mBIVdPYNlZtAjBCK8asKpH3/NyYJ8xhsBjxXLLiQifKiGOpA\n" +
            "LLBy/LyJWmo4R4zkAtUILD4FcsIyLMIJlsqWjaNdey7bwGI75hZQkBIF8QJxFVtT\n" +
            "n4HQBtuNe2ek7e72d+bayceJvlUAFXTu6oeX9/UuS7AhuY4giNzI1pNOgNwWXRxx\n" +
            "REmwvPrzJatZZ7cwfsKTezSSQlv2O4q70+2X2h0VtUg/pkz3GknE07S3ggDR9Qkg\n" +
            "bywQS/42luPIADbbAKXhHaBaX/TaD/uZVn+BOZ5sqWmxEbbHtvzlSea02J1Fk4Hq\n" +
            "kWbpuzByCJ25SuDRr+Xyn84ZDnetumQ0lBkc2ro+rZKXw8YGMyt0aX8ZwJxL4qNB\n" +
            "/WFFEproVsOru8G7iwXgt4QP8WRBSp2kTlQUbNTF3gxOTsslkUErTnvcRQ0GpK06\n" +
            "DRQG8wbjgewpHyw7O8Sfi34EjAzic0gwtIp501/MWmKpRUgAow9LPreiaLq2TBIQ\n" +
            "DXEhUb9fEhY77QKeir8cpue3sShqcz9TLa5REJGqsP/8/URk7lZjiI+YWbRLp2U2\n" +
            "D//0NPEq8fxrzNtacZRxSdx2id/yTWumtj5swjFA4yk0tunadltDMgEYuKgR+Jw9\n" +
            "G3/yFTDnepHK41V6x8eE/4JjUAvIJWADDWxudO7oF/wsY0AnUuWe9DkW09g8IWhk\n" +
            "NukDTdpsl08hCLF06qH3MSHJrdUAzs2GGLMCvtrXK2L3k70PcLqMXhbPSr7d1RGW\n" +
            "gW0BlRfR4l+2LJ952SMv3xzuxgT43aX3FFVBxXk7nFrhWJWIpJpuYXRhTqASkzoZ\n" +
            "KzsIRyW0ZbsaIsy0tgzzyhQvdoOoJn+2sKjcCzpfY6tgRD9sfucOm1sGet/cM5YP\n" +
            "iJYei2qKMeYcvACWiI8GNGY37OzhlikbleO4xXnfJwEOYx66NjTHZqkz1/TiCBGU\n" +
            "a7h+l/fnut6VfkxS1yZ2r5Gsdx7DUfNkEeKyzIMnYRA3zw3047lHqH714rV5VbE3\n" +
            "yYEQWvdtYlHMFM2z9DDta59RRATOemm7AA1fYsfodrV/QPJi5qPmvpHtCvfItbdL\n" +
            "Fg88Zh1zV5nV+0doUTXFVR9poJRE9fASlfU5qCJ9Jx5ISfvIkGz1fmfqXhUN9fE7\n" +
            "C0Evl7IYQLguTXFznRvsXvnliwR9Ut/g85JtXUiku4F2ThCBMHBDbov6p128kP+2\n" +
            "7LBgShM4IG80clxon8sWh6y0RLUz1MTamEYZKCXAPZzJoWhbzdNns/QTsjNP8wlu\n" +
            "vBRtdkb6w4Vrm6GO2BXY6pQUBPcoDuymAhfAF9TxRn860OQeMcT/NRsU9Z/8nRnz\n" +
            "3KbAuMTYsQ6qbjuLTDwfF9B4b4YUDQR22z8wlzCNLzgwFlGSI12xhf3ejRlwjGZJ\n" +
            "J/11Up4pEegRS/c+Li2OUvQr9Jxi8XGIdEJZY1T8oVpzDJf3C29gpARWSDAXrFn0\n" +
            "lgZHnqFyebeC1uDW8r/wGtYmI2EC53+FlOF5AFcH+3LzObZzerqwror4UMOA+B5c\n" +
            "QMU5vDv1LFcWLzvJHMXJfCHL5nVSukXCMawr+DbeKjrkseG0UX0gpUbQy0vHIH1K\n" +
            "2geD2xyl3TJ8jCaKOxb/Hu+KfkvtOCsh07TA+cnTV1WHR77svUcMErzHXWOFm8+U\n" +
            "omIXALO1EiDbpu38gERRLkC84eMhRBQjKcdmlcBFsmilt3cfIofypuhMRiIFjIke\n" +
            "00y2GEdQVsZGA/LX1HILqD4dEFDDQI2LPvCG5qe28HTfWspzsqK94IRESzm+Vmdp\n" +
            "IjNzkTyrPI06yMvxaHGajwUtLWCReJOG/uXhswbX7EviVYyqCR4vzDLDVXAulxo/\n" +
            "OsHaQhMX8xYOLXontx7SNCBlu/EEBww5QklKUldgd5igr7bDxsvZ6vHy/wcNIzY3\n" +
            "RUdidnuDkpSm1hIoLz4/SW2Tm6C2u9La5evu7xAfIy1ul8LE3/P0AAAAAAAAAAAA\n" +
            "AAAAABcmOEM=\n" +
            "-----END CERTIFICATE-----";

    private static String bcGeneratedMLDSA44 = "-----BEGIN PRIVATE KEY-----\n" +
            "MDQCAQAwCwYJYIZIAWUDBAMRBCIEIAABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZ\n" +
            "GhscHR4f\n" +
            "-----END PRIVATE KEY-----";

    final static String ML_DSA_44_seed = "-----BEGIN PRIVATE KEY-----\n" +
            "MDQCAQAwCwYJYIZIAWUDBAMRBCKAIAABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZ\n" +
            "GhscHR4f\n" +
            "-----END PRIVATE KEY-----";
}
