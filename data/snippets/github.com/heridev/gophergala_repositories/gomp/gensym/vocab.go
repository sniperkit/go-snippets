// The content of this file is taken from http://www.myvocabulary.com/word-list/engineering-vocabulary/
// Because what software engineering it is without words used by real engineers?
package gensym

import "strings"

const VocabOn = false

var rawVocab = `Abutment Activation Advice Advise Amplitude Analysis Angle Assembly Automation Axis
Balance Bearing Blueprint Building
Calculation Cantilever Cell Combustion Communication Component Component Compress Concept Constriction Construction Consultation Control Conversion Conveyance Cooling Coupling Crank Current
Degree Device Diagram Distill Distribution
Elastic Electrical Electronics Element Energy Engine Excavation Expert
Fabrication Flexible Flow Fluid Force Frame Fuel Fulcrum
Gear Gears Generate Generator Gimbals Grade Grading
Hardware Hoist Horizontal Hydraulic
Information Installation Instrument Intersection
Joint
Lift Load
Machine Management Manufacturing Mark Measurement Mechanize Modular Mold Motion	
Object Operation
Physics Plumb Pneumatic Power Precision Process Production Project Propulsion Pulley	
Radiate Ream Refine Regulation Repair Retrofit Rotation
Savvy Scheme Schooling Scientific Sequence Shape Slide Stability Strength Structure Structure Studying Superstructure Suspension
Technology Tools Transform Transmission Transmit Turbine	
Vacuum Valve Vertical Vibration
Weight Weld Withstand Worker`

var vocab = strings.Split(rawVocab, " ")
