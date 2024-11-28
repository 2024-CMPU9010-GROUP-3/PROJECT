const WhatIsMagpie = () => {
  return (
    <div className="py-16 px-4 max-w-6xl mx-auto" id="about">
      <h2 className="text-3xl font-bold text-center mb-8">What is Magpie?</h2>

      <div className="space-y-6 text-lg">
        <p>
          Magpie is a cutting-edge geographical information service designed
          to provide users with a comprehensive view of public services and
          amenities at a glance. It empowers urban planners, residents, and
          other stakeholders by offering interactive maps, smart search
          capabilities, and detailed analysis tools. Whether you&aposre
          planning infrastructure, exploring local services, or seeking
          information about your daily routes, Magpie simplifies access to
          essential data with reliability and precision.
        </p>
        <br />
      </div>

      <h2 className="text-3xl font-bold text-center mb-8">
        How Magpie uses Machine Learning?
      </h2>

      <div className="space-y-6 text-lg">
        <div className="space-y-6">
          <p className="text-gray-600">
            Our project uses machine learning to analyse satellite images and
            accurately identify parking spaces. Our approach relies on the
            YOLOv8 OBB Model (You Only Look Once with Oriented Bounding Boxes),
            a state-of-the-art deep learning model optimized for object
            detection.
          </p>

          <div className="space-y-4">
            <div>
              <h3 className="font-bold text-lg mb-2">
                Car Detection using the YOLO v8 OBB Model:
              </h3>
              <ul className="list-disc pl-6 space-y-2 text-gray-600">
                <li>
                  We fine-tuned the model to detect cars in satellite images.
                </li>
                <li>
                  The model was trained on our dataset of manually annotated
                  satellite images, allowing it to handle challenges commonly
                  found in satellite images such as low resolutions, shading,
                  and occlusions.
                </li>
                <li>
                  Once trained, the model outputs precise bounding boxes that
                  capture the position, size, and orientation of cars. These
                  bounding boxes are then converted to geographic coordinates.
                </li>
              </ul>
            </div>

            <div>
              <h3 className="font-bold text-lg mb-2">
                Identifying Parking Spaces:
              </h3>
              <ul className="list-disc pl-6 space-y-2 text-gray-600">
                <li>
                  We apply multiple heuristics to detect parking spaces using
                  the car detected by the model. Key steps include:
                  <ul className="list-disc pl-6 mt-2">
                    <li>
                      A road mask is applied to filter out cars on the road,
                      isolating parking lots and on the street parking.
                    </li>
                    <li>
                      Empty parking spaces are identified by analysing gaps
                      between parked cars, considering alignment, orientation,
                      and standard parking dimensions.
                    </li>
                  </ul>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default WhatIsMagpie;
