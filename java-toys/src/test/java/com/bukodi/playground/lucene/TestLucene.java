package com.bukodi.playground.lucene;

import org.apache.lucene.analysis.standard.StandardAnalyzer;
import org.apache.lucene.document.Document;
import org.apache.lucene.document.Field;
import org.apache.lucene.document.NumericDocValuesField;
import org.apache.lucene.document.TextField;
import org.apache.lucene.index.IndexWriter;
import org.apache.lucene.index.IndexWriterConfig;
import org.apache.lucene.store.Directory;
import org.apache.lucene.store.NIOFSDirectory;
import org.junit.Test;

import java.nio.file.Files;
import java.nio.file.Path;


public class TestLucene {

    @Test
    public void testInMemoryIndex() throws Exception {
        Path tmpIndexPath = Files.createTempDirectory("luceneIdx");
        Directory index = new NIOFSDirectory(tmpIndexPath);
        StandardAnalyzer analyzer = new StandardAnalyzer();
        IndexWriterConfig indexWriterConfig = new IndexWriterConfig(analyzer);
        IndexWriter writter = new IndexWriter(index, indexWriterConfig);
        Document document = new Document();

        String title = "Humpty Dumpty";
        document.add(new TextField("title", title, Field.Store.YES));
        String body = "Humpty Dumpty sat on a wall,\n" +
                "Humpty Dumpty had a great fall.\n" +
                "All the king’s horses and all the king’s men\n" +
                "Couldn’t put Humpty together again.";
        document.add(new TextField("body", body, Field.Store.YES));
        document.add( new NumericDocValuesField("added", System.currentTimeMillis()));

        writter.addDocument(document);
        writter.close();

        index.syncMetaData();

        //writter.deleteDocuments();

    }

}
