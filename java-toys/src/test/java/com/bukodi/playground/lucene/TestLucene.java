package com.bukodi.playground.lucene;

import org.apache.lucene.analysis.standard.StandardAnalyzer;
import org.apache.lucene.document.Document;
import org.apache.lucene.document.Field;
import org.apache.lucene.document.NumericDocValuesField;
import org.apache.lucene.document.TextField;
import org.apache.lucene.index.DirectoryReader;
import org.apache.lucene.index.IndexWriter;
import org.apache.lucene.index.IndexWriterConfig;
import org.apache.lucene.index.Term;
import org.apache.lucene.queries.function.valuesource.RangeMapFloatFunction;
import org.apache.lucene.search.IndexSearcher;
import org.apache.lucene.search.ScoreDoc;
import org.apache.lucene.store.ByteBuffersDirectory;
import org.apache.lucene.store.Directory;
import org.apache.lucene.util.IOUtils;
import org.junit.Test;


public class TestLucene {

    @Test
    public void testInMemoryIndex() throws Exception {
        Directory memoryIndex = new ByteBuffersDirectory();
        StandardAnalyzer analyzer = new StandardAnalyzer();
        IndexWriterConfig indexWriterConfig = new IndexWriterConfig(analyzer);
        IndexWriter writter = new IndexWriter(memoryIndex, indexWriterConfig);
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


        //writter.deleteDocuments();

    }

}
