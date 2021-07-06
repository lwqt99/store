import imageio
import os
import numpy as np
import matplotlib.pyplot as plt
from tensorflow.keras.models import Sequential,Model
from tensorflow.keras.layers import UpSampling2D, Conv2D, Dense, BatchNormalization, LeakyReLU, Input,Reshape, MaxPooling2D, Flatten, AveragePooling2D, Conv2DTranspose
from tensorflow.keras.optimizers import Adam

'''
    定义参数
'''
# 数据集相对位置
avatar_img_path = "../dcgan_anime_avatars-master/data"
# noise的维度
noise_dim = 100
# 图片的shape
image_shape = (64, 64, 3)


def load_data():
    """
    加载数据集
    :return: 返回numpy数组
    """
    all_images = []
    for image_name in os.listdir(avatar_img_path):
        # 加载图片
        image = imageio.imread(os.path.join(avatar_img_path, image_name))
        all_images.append(image)
    all_images = np.array(all_images)
    # 将图片数值变成[-1,1]
    all_images = (all_images - 127.5) / 127.5
    # 将数据随机排序
    np.random.shuffle(all_images)
    return all_images

img_dataset = load_data()

def show_images(images, index=-1):
    """
    展示并保存图片
    :param images: 需要show的图片
    :param index: 图片名
    :return:
    """
    plt.figure()
    for i, image in enumerate(images):
        ax = plt.subplot(5, 5, i + 1)
        plt.axis('off')
        plt.imshow(image)
    plt.savefig("./result/data_%d.png" % index)


def build_G():
    """
    构建生成器
    :return:
    """
    model = Sequential()
    # 全连接层 100 -> 2048
    model.add(Dense(2048, input_dim=noise_dim))
    # 激活函数
    model.add(LeakyReLU(0.2))
    # 全连接层 2048 ->  8 * 8 * 256
    model.add(Dense(8 * 8 * 256))
    # DN层
    model.add(BatchNormalization())
    model.add(LeakyReLU(0.2))
    # 8 * 8 * 256 -> (8,8,256)
    model.add(Reshape((8, 8, 256)))
    # 卷积层 (8,8,256) -> (8,8,128)
    model.add(Conv2D(128, kernel_size=5, padding='same'))
    model.add(BatchNormalization())
    model.add(LeakyReLU(0.2))
    # 反卷积层 (8,8,128) -> (16,16,128)
    model.add(Conv2DTranspose(128, kernel_size=5, strides=2, padding='same'))
    model.add(LeakyReLU(0.2))
    # 反卷积层 (16,16,128) -> (32,32,64)
    model.add(Conv2DTranspose(64, kernel_size=5, strides=2, padding='same'))
    model.add(LeakyReLU(0.2))
    # 反卷积层  (32,32,64) -> (64,64,3) = 图片
    model.add(Conv2DTranspose(3, kernel_size=5, strides=2, padding='same', activation='tanh'))
    return model


def build_D():
    """
    构建判别器
    :return:
    """
    model = Sequential()
    # 卷积层
    model.add(Conv2D(64, kernel_size=5, padding='valid', input_shape=image_shape))
    # BN层
    model.add(BatchNormalization())
    # 激活层
    model.add(LeakyReLU(0.2))
    # 平均池化层
    model.add(AveragePooling2D(pool_size=2))
    # 卷积层
    model.add(Conv2D(128, kernel_size=3, padding='valid'))
    model.add(BatchNormalization())
    model.add(LeakyReLU(0.2))
    model.add(AveragePooling2D(pool_size=2))
    model.add(Conv2D(256, kernel_size=3, padding='valid'))
    model.add(BatchNormalization())
    model.add(LeakyReLU(0.2))
    model.add(AveragePooling2D(pool_size=2))
    # 将输入展平
    model.add(Flatten())
    # 全连接层
    model.add(Dense(1024))
    model.add(BatchNormalization())
    model.add(LeakyReLU(0.2))
    # 最终输出1(true img) 0(fake img)的概率大小
    model.add(Dense(1, activation='sigmoid'))
    model.compile(loss='binary_crossentropy',
                  optimizer=Adam(learning_rate=0.0002, beta_1=0.5))
    return model


def build_gan():
    """
    构建GAN网络
    :return:
    """
    # 冷冻判别器，也就是在训练的时候只优化G的网络权重，而对D保持不变
    D.trainable = False
    # GAN网络的输入
    gan_input = Input(shape=(noise_dim,))
    # GAN网络的输出
    gan_out = D(G(gan_input))
    # 构建网络
    gan = Model(gan_input, gan_out)
    # 编译GAN网络，使用Adam优化器，以及加上交叉熵损失函数（一般用于二分类）
    gan.compile(loss='binary_crossentropy', optimizer=Adam(learning_rate=0.0002, beta_1=0.5))
    return gan


def sample_noise(batch_size):
    """
    随机产生正态分布（0，1）的noise
    :param batch_size:
    :return: 返回的shape为(batch_size,noise)
    """
    return np.random.normal(size=(batch_size, noise_dim))


def load_batch(data, batch_size, index):
    """
    按批次加载图片
    :param data: 图片数据集
    :param batch_size: 批次大小
    :param index: 批次序号
    :return:
    """
    return data[index * batch_size: (index + 1) * batch_size]


def train(epochs=100, batch_size=64):
    """
    训练函数
    :param epochs: 训练的次数
    :param batch_size: 批尺寸
    :return:
    """
    # 判别器损失
    discriminator_loss = 0
    # 生成器损失
    generator_loss = 0
    # img_dataset.shape[0] / batch_size 代表这个数据可以分为几个批次进行训练
    n_batches = int(img_dataset.shape[0] / batch_size)

    for i in range(epochs):
        for index in range(n_batches):
            # 按批次加载数据
            x = load_batch(img_dataset, batch_size, index)
            # 产生noise
            noise = sample_noise(batch_size)
            # G网络产生图片
            generated_images = G.predict(noise)
            #show_images(generated_images[0:25], i*n_batches+index)
            # 产生为1的标签
            y_real = np.ones(batch_size)
            # 产生为0的标签
            y_fake = np.zeros(batch_size)
            # 训练真图片loss
            d_loss_real = D.train_on_batch(x, y_real)
            # 训练假图片loss
            d_loss_fake = D.train_on_batch(generated_images, y_fake)
            discriminator_loss = d_loss_real + d_loss_fake

            noise = sample_noise(batch_size)
            # 训练GAN网络，input = fake_img ,label = 1
            generator_loss = GAN.train_on_batch(noise, y_real)

        if i % 10 == 0:
            if i == 0:
                pass
            else:
                G.save_weights("./models/Gene"+str(i)+".hdf5")

        # 随机产生(25,100)的noise
        test_noise = sample_noise(25)
        # 使用G网络生成25张图片
        test_images = G.predict(test_noise)
        show_images(test_images, i)
        print(
            '[Epoch {0}]. Discriminator loss : {1}. Generator_loss: {2}.'.format(i, discriminator_loss, generator_loss))


G = build_G()
D = build_D()
GAN = build_gan()
train(epochs=500, batch_size=128)
