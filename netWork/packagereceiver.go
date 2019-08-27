package main

improt(
	"projects/encodeImages/models"
)

var(
	dataSizePerFrame = 1600 * 10
    frameBuffer []byte
    oneFrameBuffer []byte
    cursize = 0
)

    func Init(){
         pool_size := 2 * BigScreen.width * BigScreen.height * 3 / dataSizePerFrame
    }

    //接收服务器数据
     func ReciveMsg(){
        buffer_size := totalSizePerFrame * 5;
        byte[] buffer = new byte[buffer_size];
        byte[] onepack = new byte[totalSizePerFrame];
        List<PackData> packlist = new List<PackData>();
        while (true)
        {
            int length = client.ReceiveFrom(buffer, ref point);//接收数据报

            //ReadFrameData(buffer, length);

            packlist.Clear();
            int curpos = 0;

            while ((curpos + totalSizePerFrame) <= buffer_size && (curpos + totalSizePerFrame) <= length)
            {
                Array.Copy(buffer, curpos, onepack, 0, totalSizePerFrame);
                packlist.Add(ReadOnePack(onepack));
                curpos += totalSizePerFrame;
            }
            JoinPackData(packlist);
        }
    }

    //read one pack
    func ReadOnePack(data []byte ) PackData{
        PackData pack = packDataPool.Create();
        pack.datasize = System.BitConverter.ToInt32(data, 0);
        pack.cursize = System.BitConverter.ToInt32(data, 4);
        pack.last_pak = System.BitConverter.ToInt32(data, 8);
        Array.Copy(data, 12, pack.data, 0, dataSizePerFrame);
        return pack;
    }

    //joint pack list
    func JoinPackData(list []PackData){
        for  i := 0; i < list.Count; i++{
             data := list[i];
            if frameBuffer == null || frameBuffer.Length != data.datasize{
                frameBuffer = make([]byte,data.datasize);
            }
            //若为最後一片
            if (data.last_pak == 1){
                //成功获取一个数据包
                if ((cursize + data.cursize) == data.datasize){
                    Array.Copy(data.data, 0, frameBuffer, cursize, data.cursize);
                    if (oneFrameBuffer == null || oneFrameBuffer.Length != data.datasize) {
                        oneFrameBuffer = new byte[data.datasize];
                    }

                    lock (oneFrameBuffer){
                        frameBuffer.CopyTo(oneFrameBuffer, 0);
                        cursize = 0;
                        CreateColors(oneFrameBuffer);
                        BigScreen.flag = true;
                    }
                }else{
                    //丢弃该帧
                    cursize = 0;
                }
            }
            else if cursize < data.datasize && (cursize + data.cursize) <= data.datasize{
                //往下拼接
                data.data.CopyTo(frameBuffer, cursize);
                cursize += data.cursize;
            }else{
                //丢弃该帧
                cursize = 0;
            }
        }
    }

     func GetCurFrameBuffer() []byte{
        return oneFrameBuffer;
    }